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

// Client struct
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

func (c *Client) call(method string, path string, body any) (io.ReadCloser, error) {
	ep := "https://api.dribbble.com/v2" + path
	u, err := url.Parse(ep)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
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

	dribResp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if dribResp != nil && c.Flags.Verbose {
		grey := color.New(color.FgHiBlack).SprintFunc()
		s := fmt.Sprintf("RateLimit-> %d, RateLimitRemaining-> %d, RateLimitReset-> %s",
			dribResp.RateLimit,
			dribResp.RateLimitRemaining,
			time.Unix(dribResp.RateLimitReset, 0).Format("01/02/2006 03:04 PM MST"),
		)
		fmt.Println(grey(s))
	}

	return dribResp.Body, err
}

type DribbbleResponse struct {
	Body               io.ReadCloser
	StatusCode         int
	RateLimit          int
	RateLimitRemaining int
	RateLimitReset     int64
}

func (c *Client) do(req *http.Request) (*DribbbleResponse, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	dribResp := &DribbbleResponse{
		Body:       res.Body,
		StatusCode: res.StatusCode,
	}

	// extract rate limit headers
	dribResp.RateLimit, _ = strconv.Atoi(res.Header.Get("X-RateLimit-Limit"))
	dribResp.RateLimitRemaining, _ = strconv.Atoi(res.Header.Get("X-RateLimit-Remaining"))
	if reset, _ := strconv.ParseInt(res.Header.Get("X-RateLimit-Reset"), 10, 64); reset != 0 {
		dribResp.RateLimitReset = reset
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
