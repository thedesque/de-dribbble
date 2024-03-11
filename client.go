package dribbble

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Client for Dribbble API.
// See https://developer.dribbble.com/v2/
type Client struct {
	*Config
	User        *User
	Projects    *Projects
	Shots       *Shots
	Jobs        *Jobs
	Likes       *Likes
	Attachments *Attachments
}

// NewClient returns new instance of Dribbble client
func NewClient(config *Config) *Client {
	c := &Client{Config: config}
	c.User = &User{c}
	c.Projects = &Projects{c}
	c.Shots = &Shots{c}
	c.Jobs = &Jobs{c}
	c.Likes = &Likes{c}
	c.Attachments = &Attachments{c}
	return c
}

// ------------------------------------------------------------------------

// client.call performs the HTTP request and returns the response body.
func (c *Client) call(method string, path string, body any) (*dribbbleResponse, error) {
	ep := "https://api.dribbble.com/v2" + path
	u, err := url.Parse(ep)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	// use body argument if it was specified
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if res != nil && c.Flags.Verbose {
		s := fmt.Sprintf("RateLimit-> %d, RateLimitRemaining-> %d, RateLimitReset-> %s",
			res.rateLimit,
			res.rateLimitRemaining,
			time.Unix(res.rateLimitReset, 0).Format("01/02/2006 03:04 PM MST"),
		)
		fmt.Println(color.HiBlackString(s))
	}

	return res, err
}

// dribbbleResponse represents a response from the Dribbble API.
// Stores the response body, status code, rate limit, and pagination info.
type dribbbleResponse struct {
	body               io.ReadCloser
	statusCode         int
	rateLimit          int
	rateLimitRemaining int
	rateLimitReset     int64
	pagination         pagination
}

// client.do performs the HTTP request and returns the response.
func (c *Client) do(req *http.Request) (*dribbbleResponse, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	dribResp := &dribbbleResponse{
		body:       res.Body,
		statusCode: res.StatusCode,
	}

	// parse link header for pagination URLs
	if linkHeader := res.Header.Get("Link"); linkHeader != "" {
		paginationInfo := parseLinkHeader(linkHeader)
		dribResp.pagination = paginationInfo
	}

	// extract rate limit headers
	dribResp.rateLimit, _ = strconv.Atoi(res.Header.Get("X-RateLimit-Limit"))
	dribResp.rateLimitRemaining, _ = strconv.Atoi(res.Header.Get("X-RateLimit-Remaining"))
	if reset, _ := strconv.ParseInt(res.Header.Get("X-RateLimit-Reset"), 10, 64); reset != 0 {
		dribResp.rateLimitReset = reset
	}

	if res.StatusCode < 400 {
		return dribResp, nil
	}

	defer res.Body.Close()

	e := &Error{
		StatusCode: res.StatusCode,
		Message:    res.Status,
	}

	ct := res.Header.Get("Content-Type")
	if strings.Contains(ct, "text/html") {
		return nil, e
	}

	if err := json.NewDecoder(res.Body).Decode(e); err != nil {
		return nil, err
	}

	return nil, e
}

// ------------------------------------------------------------------------

// pagination represents pagination information for a Dribbble API response.
// It includes the URLs for the next and previous pages, the number of items per page,
// and the previous, current, and next page numbers.
type pagination struct {
	nextPageURL string
	prevPageURL string
	perPage     int
	prevPage    int
	currentPage int
	nextPage    int
}

// parseLinkHeader extracts pagination information from a Link header provided by the Dribbble API.
// It parses the header to retrieve URLs for the next and previous pages, along with pagination
// parameters such as the number of items per page (perPage), the previous page number (prevPage),
// the current page number (currentPage), and the next page number (nextPage). The function assumes
// that if a "next" page link is present, the current page is one less than the "next" page; if a
// "prev" page link is present and a "next" page link is not, the current page is one more than
// the "prev" page.
func parseLinkHeader(header string) pagination {
	var pageInfo pagination
	links := strings.Split(header, ",")
	for _, link := range links {
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) == 2 {
			urlPart := strings.TrimSpace(parts[0])
			urlPart = strings.Trim(urlPart, "<>")
			relPart := strings.TrimSpace(parts[1])

			u, err := url.Parse(urlPart)
			if err != nil {
				continue
			}

			queryParams := u.Query()
			page, err := strconv.Atoi(queryParams.Get("page"))
			if err != nil {
				continue
			}

			perPage, err := strconv.Atoi(queryParams.Get("per_page"))
			if err != nil {
				continue
			}

			if strings.Contains(relPart, `rel="next"`) {
				pageInfo.nextPageURL = urlPart
				pageInfo.nextPage = page
				pageInfo.perPage = perPage
				// currentPage is one less than nextPage since nextPage exists
				pageInfo.currentPage = page - 1
			} else if strings.Contains(relPart, `rel="prev"`) {
				pageInfo.prevPageURL = urlPart
				pageInfo.prevPage = page
				pageInfo.perPage = perPage
				// currentPage is one more than prevPage since prevPage exists
				pageInfo.currentPage = page + 1
			}
		}
	}

	// adjusts currentPage based on nextPage and prevPage if necessary
	if pageInfo.currentPage == 0 {
		if pageInfo.nextPage > 0 {
			pageInfo.currentPage = pageInfo.nextPage - 1
		} else if pageInfo.prevPage > 0 {
			pageInfo.currentPage = pageInfo.prevPage + 1
		}
	}

	return pageInfo
}
