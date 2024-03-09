package dribbble

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// User client
type User struct {
	*Client
}

// UserOut defines the structure of user information.
type UserOut struct {
	ID        int    `json:"id,omitempty" toml:"id, omitempty"`
	Name      string `json:"name,omitempty" toml:"name,omitempty"`
	Login     string `json:"login,omitempty" toml:"login,omitempty"`
	HTMLURL   string `json:"html_url,omitempty" toml:"htmlurl,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty" toml:"avatar_url,omitempty"`
	Bio       string `json:"bio,omitempty" toml:"bio,omitempty"`
	Location  string `json:"location,omitempty" toml:"location,omitempty"`
	Links     struct {
		Web     string `json:"web,omitempty" toml:"web,omitempty"`
		Twitter string `json:"twitter,omitempty" toml:"twitter,omitempty"`
	} `json:"links,omitempty" toml:"links,omitempty"`
	CanUploadShot  bool      `json:"can_upload_shot,omitempty" toml:"can_upload_shot"`
	Pro            bool      `json:"pro,omitempty" toml:"pro"`
	FollowersCount int       `json:"followers_count,omitempty" toml:"followers_count,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty" toml:"created_at"`
	Type           string    `json:"type,omitempty" toml:"type,omitempty"`
	Teams          []struct {
		ID        int    `json:"id,omitempty" toml:"id,omitempty"`
		Name      string `json:"name,omitempty" toml:"name,omitempty"`
		Login     string `json:"login,omitempty" toml:"login,omitempty"`
		HTMLURL   string `json:"html_url,omitempty" toml:"htmlurl,omitempty"`
		AvatarURL string `json:"avatar_url,omitempty" toml:"avatar_url,omitempty"`
		Bio       string `json:"bio,omitempty" toml:"bio,omitempty"`
		Location  string `json:"location,omitempty" toml:"location,omitempty"`
		Links     struct {
			Web     string `json:"web,omitempty" toml:"web,omitempty"`
			Twitter string `json:"twitter,omitempty" toml:"twitter,omitempty"`
		} `json:"links,omitempty" toml:"links, omitempty"`
		Type      string    `json:"type,omitempty" toml:"type,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty" toml:"created_at"`
		UpdatedAt time.Time `json:"updated_at,omitempty" toml:"updated_at"`
	} `json:"teams,omitempty" toml:"teams,omitempty"`
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

	writeTitleString(&sb, "User")
	writeIfNotEmpty(&sb, "User", fmt.Sprintf("%s (%s)", u.Name, u.Login))
	writeIfNotEmpty(&sb, "Profile", u.HTMLURL)
	writeIfNotEmpty(&sb, "Bio", u.Bio)
	writeIfNotEmpty(&sb, "Location", u.Location)
	writeIfNotEmpty(&sb, "Web:", u.Links.Web)
	writeIfNotEmpty(&sb, "Twitter:", u.Links.Twitter)
	writeIfNotEmpty(&sb, "Pro", fmt.Sprintf("%t", u.Pro))
	writeIfNotEmpty(&sb, "Can Upload Shot", fmt.Sprintf("%t", u.CanUploadShot))
	writeIfNotEmpty(&sb, "Followers", fmt.Sprintf("%d", u.FollowersCount))
	writeIfNotEmpty(&sb, "Created At", u.CreatedAt.Format("Jan 2, 2006"))

	if len(u.Teams) > 0 {
		writeTitleString(&sb, "Teams")
		for _, team := range u.Teams {
			teamDetails := fmt.Sprintf("%s (%s): %s", team.Name, team.Type, team.Bio)
			writeIfNotEmpty(&sb, "-", teamDetails)
		}
	} else {
		writeIfNotEmpty(&sb, "Teams", "None")
	}

	return sb.String()
}

func (out *UserOut) ToToml() (string, error) {
	return toTomlString(out)
}
