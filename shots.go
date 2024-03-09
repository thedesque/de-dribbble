package dribbble

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Shots instance
type Shots struct {
	*Client
}

// ShotOut single schema
type ShotOut struct {
	ID          int    `json:"id,omitempty" toml:"id,omitempty"`
	Title       string `json:"title,omitempty" toml:"title,omitempty"`
	Description string `json:"description,omitempty" toml:"description,omitempty"`
	Width       int    `json:"width,omitempty" toml:"width,omitempty"`
	Height      int    `json:"height,omitempty" toml:"height,omitempty"`
	Images      struct {
		Hidpi  any    `json:"hidpi,omitempty" toml:"hidpi,omitempty"`
		Normal string `json:"normal,omitempty" toml:"normal,omitempty"`
		Teaser string `json:"teaser,omitempty" toml:"teaser,omitempty"`
	} `json:"images,omitempty" toml:"images,omitempty"`
	PublishedAt time.Time `json:"published_at,omitempty" toml:"published_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" toml:"updated_at"`
	HTMLURL     string    `json:"html_url,omitempty" toml:"html_url,omitempty"`
	Animated    bool      `json:"animated,omitempty" toml:"animated"`
	Tags        []string  `json:"tags,omitempty" toml:"tags,omitempty"`
	Attachments []struct {
		ID           int       `json:"id,omitempty" toml:"id,omitempty"`
		URL          string    `json:"url,omitempty" toml:"url,omitempty"`
		ThumbnailURL string    `json:"thumbnail_url,omitempty" toml:"thumbnail_url,omitempty"`
		Size         int       `json:"size,omitempty" toml:"size,omitempty"`
		ContentType  string    `json:"content_type,omitempty" toml:"content_type,omitempty"`
		CreatedAt    time.Time `json:"created_at,omitempty" toml:"created_at"`
	} `json:"attachments,omitempty" toml:"attachments,omitempty"`
	Projects []struct {
		ID          int       `json:"id,omitempty" toml:"id,omitempty"`
		Name        string    `json:"name,omitempty" toml:"name,omitempty"`
		Description string    `json:"description,omitempty" toml:"description,omitempty"`
		ShotsCount  int       `json:"shots_count,omitempty" toml:"shots_count,omitempty"`
		CreatedAt   time.Time `json:"created_at,omitempty" toml:"created_at"`
		UpdatedAt   time.Time `json:"updated_at,omitempty" toml:"updated_at"`
	} `json:"projects,omitempty" toml:"projects,omitempty"`
	Team struct {
		ID        int    `json:"id,omitempty" toml:"id,omitempty"`
		Name      string `json:"name,omitempty" toml:"name,omitempty"`
		Login     string `json:"login,omitempty" toml:"login,omitempty"`
		HTMLURL   string `json:"html_url,omitempty" toml:"html_url,omitempty"`
		AvatarURL string `json:"avatar_url,omitempty" toml:"avatar_url,omitempty"`
		Bio       string `json:"bio,omitempty" toml:"bio,omitempty"`
		Location  string `json:"location,omitempty" toml:"location,omitempty"`
		Links     struct {
			Web     string `json:"web,omitempty" toml:"web,omitempty"`
			Twitter string `json:"twitter,omitempty" toml:"twitter,omitempty"`
		} `json:"links,omitempty" toml:"links,omitempty"`
		Type      string    `json:"type,omitempty" toml:"type,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty" toml:"created_at"`
		UpdatedAt time.Time `json:"updated_at,omitempty" toml:"updated_at"`
	} `json:"team,omitempty" toml:"team,omitempty"`
	Video struct {
		ID               int       `json:"id,omitempty" toml:"id,omitempty"`
		Duration         int       `json:"duration,omitempty" toml:"duration,omitempty"`
		VideoFileName    string    `json:"video_file_name,omitempty" toml:"video_file_name,omitempty"`
		VideoFileSize    int       `json:"video_file_size,omitempty" toml:"video_file_size,omitempty"`
		Width            int       `json:"width,omitempty" toml:"width,omitempty"`
		Height           int       `json:"height,omitempty" toml:"height,omitempty"`
		Silent           bool      `json:"silent,omitempty" toml:"silent"`
		CreatedAt        time.Time `json:"created_at,omitempty" toml:"created_at"`
		UpdatedAt        time.Time `json:"updated_at,omitempty" toml:"updated_at"`
		URL              string    `json:"url,omitempty" toml:"url,omitempty"`
		SmallPreviewURL  string    `json:"small_preview_url,omitempty" toml:"small_preview_url,omitempty"`
		LargePreviewURL  string    `json:"large_preview_url,omitempty" toml:"large_preview_url,omitempty"`
		XlargePreviewURL string    `json:"xlarge_preview_url,omitempty" toml:"xlarge_preview_url,omitempty"`
	} `json:"video,omitempty" toml:"video,omitempty"`
	LowProfile bool `json:"low_profile,omitempty" toml:"low_profile"`
}

// PopularShotOut schema
type PopularShotOut struct {
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
}

// UpdateShotIn for updating shot
type UpdateShotIn struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// ------------------------------------------------------------------------

// GetShots of authenticated user
func (c *Shots) GetShots() (out *[]ShotOut, err error) {
	body, err := c.call("GET", "/user/shots", nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// String method to convert ShotOut struct into a human-readable string,
// with colored keys and omitting empty values.
func (s *ShotOut) String() string {
	var sb strings.Builder

	writeTitleString(&sb, "Shot")
	writeIfNotEmpty(&sb, "ID", fmt.Sprintf("%d", s.ID))
	writeIfNotEmpty(&sb, "Title", s.Title)
	writeIfNotEmpty(&sb, "Description", s.Description)
	writeIfNotEmpty(&sb, "Dimensions", fmt.Sprintf("%dx%d", s.Width, s.Height))
	writeIfNotEmpty(&sb, "HTML URL", s.HTMLURL)
	writeIfNotEmpty(&sb, "Image - HiDPI:", s.Images.Hidpi.(string))
	writeIfNotEmpty(&sb, "Image - Normal:", s.Images.Normal)
	writeIfNotEmpty(&sb, "Image - Teaser:", s.Images.Teaser)
	writeIfNotEmpty(&sb, "Published At", s.PublishedAt.Format("Jan 2, 2006"))
	writeIfNotEmpty(&sb, "Updated At", s.UpdatedAt.Format("Jan 2, 2006"))
	writeIfNotEmpty(&sb, "Animated", fmt.Sprintf("%t", s.Animated))
	writeIfNotEmpty(&sb, "Tags:", formatTags(s.Tags))

	if len(s.Attachments) > 0 {
		writeTitleString(&sb, "Attachments")
		for _, attachment := range s.Attachments {
			attachmentDetails := fmt.Sprintf("ID: %d, URL: %s, Size: %d bytes", attachment.ID, attachment.URL, attachment.Size)
			writeIfNotEmpty(&sb, "-", attachmentDetails)
		}
	}

	if len(s.Projects) > 0 {
		writeTitleString(&sb, "Projects")
		for _, project := range s.Projects {
			projectDetails := fmt.Sprintf("Name: %s, Description: %s, Shots Count: %d", project.Name, project.Description, project.ShotsCount)
			writeIfNotEmpty(&sb, "-", projectDetails)
		}
	}

	if s.Team.ID != 0 {
		writeTitleString(&sb, "Team")
		teamDetails := fmt.Sprintf("Name: %s (%s), Bio: %s", s.Team.Name, s.Team.Login, s.Team.Bio)
		writeIfNotEmpty(&sb, "-", teamDetails)
	}

	if s.Video.ID != 0 {
		writeTitleString(&sb, "Video")
		videoDetails := fmt.Sprintf("Duration: %d seconds, URL: %s", s.Video.Duration, s.Video.URL)
		writeIfNotEmpty(&sb, "-", videoDetails)
	}

	writeIfNotEmpty(&sb, "Low Profile", fmt.Sprintf("%t", s.LowProfile))

	return sb.String()
}

func (out *ShotOut) ToToml() (string, error) {
	return toTomlString(out)
}

// ------------------------------------------------------------------------

// GetShot with given id
// This method returns only shots owned by the currently authenticated user
func (c *Shots) GetShot(id int) (out *ShotOut, err error) {
	body, err := c.call("GET", fmt.Sprintf("/shots/%d", id), nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// ------------------------------------------------------------------------

// GetPopularShots overall
// Note: This is available only to select applications with dribbble approval
func (c *Shots) GetPopularShots() (out *[]PopularShotOut, err error) {
	body, err := c.call("GET", "/popular_shots", nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

func (s *PopularShotOut) String() string {
	var sb strings.Builder

	writeTitleString(&sb, "Popular Shot")
	writeIfNotEmpty(&sb, "ID", fmt.Sprintf("%d", s.ID))
	writeIfNotEmpty(&sb, "Title", s.Title)
	writeIfNotEmpty(&sb, "Description", s.Description)
	writeIfNotEmpty(&sb, "Dimensions", fmt.Sprintf("%dx%d", s.Width, s.Height))
	writeIfNotEmpty(&sb, "HTML URL", s.HTMLURL)
	writeIfNotEmpty(&sb, "Image - HiDPI:", s.Images.Hidpi.(string))
	writeIfNotEmpty(&sb, "Image - Normal:", s.Images.Normal)
	writeIfNotEmpty(&sb, "Image - Teaser:", s.Images.Teaser)
	writeIfNotEmpty(&sb, "Published At", s.PublishedAt.Format("Jan 2, 2006"))

	return sb.String()
}

func (out *PopularShotOut) ToToml() (string, error) {
	return toTomlString(out)
}

// ------------------------------------------------------------------------

// UpdateShot with given id and payload
// Updating a shot requires the user to be authenticated with the upload scope
// The authenticated user must also own the shot
func (c *Shots) UpdateShot(id int, in *UpdateShotIn) (out *ShotOut, err error) {
	body, err := c.call("PUT", fmt.Sprintf("/shots/%d", id), in)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// DeleteShot with given id
// Deleting a shot requires the user to be authenticated with the upload scope
// The authenticated user must also own the shot
func (c *Shots) DeleteShot(id int) error {
	body, err := c.call("DELETE", fmt.Sprintf("/shots/%d", id), nil)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}
