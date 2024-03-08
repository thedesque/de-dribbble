package dribbble

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

// User client
type User struct {
	*Client
}

// UserOut defines the structure of user information.
type UserOut struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	HTMLURL   string `json:"html_url"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Links     struct {
		Web     string `json:"web"`
		Twitter string `json:"twitter"`
	} `json:"links"`
	CanUploadShot  bool      `json:"can_upload_shot"`
	Pro            bool      `json:"pro"`
	FollowersCount int       `json:"followers_count"`
	CreatedAt      time.Time `json:"created_at"`
	Type           string    `json:"type"`
	Teams          []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Login     string `json:"login"`
		HTMLURL   string `json:"html_url"`
		AvatarURL string `json:"avatar_url"`
		Bio       string `json:"bio"`
		Location  string `json:"location"`
		Links     struct {
			Web     string `json:"web"`
			Twitter string `json:"twitter"`
		} `json:"links"`
		Type      string    `json:"type"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"teams"`
}

// GetUser which is currenlty logged in
func (c *User) GetUser() (out *UserOut, err error) {
	body, err := c.call("GET", "/user", nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// String method to convert UserOut struct into a human-readable string,
// with colored keys and omitting empty values.
func (u UserOut) String() string {
	var sb strings.Builder
	grey := color.New(color.FgHiBlack).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	writeIfNotEmpty := func(key, value string) {
		if value != "" {
			sb.WriteString(fmt.Sprintf("%s: %s\n", grey(key), value))
		}
	}

	sb.WriteString(green("User Details:\n"))
	writeIfNotEmpty("User", fmt.Sprintf("%s (%s)", u.Name, u.Login))
	writeIfNotEmpty("Profile", u.HTMLURL)
	writeIfNotEmpty("Bio", u.Bio)
	writeIfNotEmpty("Location", u.Location)
	if u.Links.Web != "" || u.Links.Twitter != "" {
		writeIfNotEmpty("Web", u.Links.Web)
		writeIfNotEmpty("Twitter", u.Links.Twitter)
	}
	writeIfNotEmpty("Pro", fmt.Sprintf("%t", u.Pro))
	writeIfNotEmpty("Can Upload Shot", fmt.Sprintf("%t", u.CanUploadShot))
	writeIfNotEmpty("Followers", fmt.Sprintf("%d", u.FollowersCount))
	writeIfNotEmpty("Created At", u.CreatedAt.Format("Jan 2, 2006"))

	if len(u.Teams) > 0 {
		sb.WriteString(fmt.Sprintf("%s:\n", grey("Teams")))
		for _, team := range u.Teams {
			teamDetails := fmt.Sprintf("%s (%s): %s", team.Name, team.Type, team.Bio)
			writeIfNotEmpty("-", teamDetails)
		}
	} else {
		writeIfNotEmpty("Teams", "None")
	}

	return sb.String()
}
