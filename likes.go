package dribbble

import (
	"encoding/json"
	"fmt"
	"time"
)

// Likes client
type Likes struct {
	*Client
}

// LikeOut response structure
type LikeOut struct {
	ID        int       `json:"id,omitempty" toml:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" toml:"created_at"`
	Shot      struct {
		ID          int    `json:"id,omitempty" toml:"id,omitempty"`
		Title       string `json:"title,omitempty" toml:"title,omitempty"`
		Description string `json:"description,omitempty" toml:"description,omitempty"`
		Images      struct {
			Hidpi  any    `json:"hidpi,omitempty" toml:"hidpi,omitempty"`
			Normal string `json:"normal,omitempty" toml:"normal,omitempty"`
			Teaser string `json:"teaser,omitempty" toml:"teaser,omitempty"`
		} `json:"images,omitempty" toml:"images,omitempty"`
		PublishedAt time.Time `json:"published_at,omitempty" toml:"published_at"`
		HTMLURL     string    `json:"html_url,omitempty" toml:"htmlurl,omitempty"`
		Height      int       `json:"height,omitempty" toml:"height,omitempty"`
		Width       int       `json:"width,omitempty" toml:"width,omitempty"`
	} `json:"shot,omitempty" toml:"shot"`
	User struct {
		ID      int    `json:"id,omitempty" toml:"id,omitempty"`
		Name    string `json:"name,omitempty" toml:"name,omitempty"`
		Login   string `json:"login,omitempty" toml:"login,omitempty"`
		HTMLURL string `json:"html_url,omitempty" toml:"htmlurl,omitempty"`
	} `json:"user,omitempty" toml:"user,omitempty"`
}

// LikedShotOut response structure
type LikedShotOut struct {
	ID        int       `json:"id,omitempty" toml:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" toml:"created_at"`
}

// GetLikes returns list of authenticated user’s liked shots
// Note: This is available only to select applications with dribbble approval
func (c *Likes) GetLikes() (out *[]LikeOut, err error) {
	body, err := c.call("GET", "/user/likes", nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// GetShotLike checks if you like a shot
// Note: This is available only to select applications with dribbble approval
func (c *Likes) GetShotLike(id int) (out *LikedShotOut, err error) {
	body, err := c.call("GET", fmt.Sprintf("/shots/%d/like", id), nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// LikeShot with given id
// Note: This is available only to select applications with dribbble approval
func (c *Likes) LikeShot(id int) (out *LikedShotOut, err error) {
	body, err := c.call("POST", fmt.Sprintf("/shots/%d/like", id), nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// UnlikeShot with given id
// Note: This is available only to select applications with dribbble approval
// Unliking a shot requires the user to be authenticated with the write scope
func (c *Likes) UnlikeShot(id int) error {
	body, err := c.call("DELETE", fmt.Sprintf("/shots/%d/like", id), nil)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}
