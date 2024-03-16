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

// GetProjects of authenticated user. Set page to 1 to get the first page of results.
// If traverse is true, it will traverse all pages and return all projects.
func (c *Projects) GetProjects(page int, traverse bool) ([]*ProjectOut, error) {
	var projects []*ProjectOut

	for {
		path := "/user/projects" + paginationQueryString(page, 100)

		resp, err := c.call("GET", path, nil)
		if err != nil {
			return nil, err
		}
		defer resp.body.Close()

		var pageProjects []*ProjectOut
		err = json.NewDecoder(resp.body).Decode(&pageProjects)
		if err != nil {
			return nil, err
		}

		projects = append(projects, pageProjects...)

		// check if we need to traverse and if there is a next page
		if !traverse || resp.pagination.nextPage == 0 {
			break
		}

		// set up for the next iteration
		page = resp.pagination.nextPage
	}

	return projects, nil
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

// ------------------------------------------------------------------------

// CreateProject with given payload
func (c *Projects) CreateProject(in *ProjectIn) (out *ProjectOut, err error) {
	resp, err := c.call("POST", "/projects", in)
	if err != nil {
		return nil, err
	}
	defer resp.body.Close()

	err = json.NewDecoder(resp.body).Decode(&out)
	return
}

// UpdateProject with given id
func (c *Projects) UpdateProject(id int, in *ProjectIn) (out *ProjectOut, err error) {
	resp, err := c.call("PUT", fmt.Sprintf("/projects/%d", id), in)
	if err != nil {
		return nil, err
	}
	defer resp.body.Close()

	err = json.NewDecoder(resp.body).Decode(&out)
	return
}

// DeleteProject with given id
func (c *Projects) DeleteProject(id int) error {
	resp, err := c.call("DELETE", fmt.Sprintf("/projects/%d", id), nil)
	if err != nil {
		return err
	}
	defer resp.body.Close()

	return nil
}
