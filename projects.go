package dribbble

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Projects client
type Projects struct {
	*Client
}

// ProjectOut response structure
type ProjectOut struct {
	ID          int       `json:"id,omitempty" toml:"id,omitempty" yaml:"id,omitempty"`
	Name        string    `json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty"`
	Description string    `json:"description,omitempty" toml:"description,omitempty" yaml:"description,omitempty"`
	ShotsCount  int       `json:"shots_count,omitempty" toml:"shots_count,omitempty" yaml:"shots_count,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

// ProjectIn payload structure
type ProjectIn struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ------------------------------------------------------------------------

// GetProjects of authenticated user
func (c *Projects) GetProjects() (out *[]ProjectOut, err error) {
	body, err := c.call("GET", "/user/projects", nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

func (s *ProjectOut) String() string {
	var sb strings.Builder

	writeTitleString(&sb, "Project")
	writeIfNotEmpty(&sb, "ID", fmt.Sprintf("%d", s.ID))
	writeIfNotEmpty(&sb, "Name", s.Name)
	writeIfNotEmpty(&sb, "Description", s.Description)
	writeIfNotEmpty(&sb, "Shots Count", fmt.Sprintf("%d", s.ShotsCount))
	writeIfNotEmpty(&sb, "Created At", s.CreatedAt.Format("Jan 2, 2006"))
	writeIfNotEmpty(&sb, "Updated At", s.UpdatedAt.Format("Jan 2, 2006"))

	return sb.String()
}

func (out *ProjectOut) ToToml() (string, error) {
	return toTomlString(out)
}

func (out *ProjectOut) ToYaml() (string, error) {
	return toYamlString(out)
}

// ------------------------------------------------------------------------

// CreateProject with given payload
func (c *Projects) CreateProject(in *ProjectIn) (out *ProjectOut, err error) {
	body, err := c.call("POST", "/projects", in)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// UpdateProject with given id
func (c *Projects) UpdateProject(id int, in *ProjectIn) (out *ProjectOut, err error) {
	body, err := c.call("PUT", fmt.Sprintf("/projects/%d", id), in)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// DeleteProject with given id
func (c *Projects) DeleteProject(id int) error {
	body, err := c.call("DELETE", fmt.Sprintf("/projects/%d", id), nil)
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}
