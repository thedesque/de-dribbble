package dribbble

import (
	"encoding/json"
	"fmt"
	"time"
)

// Jobs client
// In order to use this part of API, you will need special token
// https://developer.dribbble.com/v2/jobs/
type Jobs struct {
	*Client
}

// JobOut schema
type JobOut struct {
	ID               int       `json:"id,omitempty" toml:"id,omitempty"`
	OrganizationName string    `json:"organization_name,omitempty" toml:"organization_name,omitempty"`
	Title            string    `json:"title,omitempty" toml:"title,omitempty"`
	Location         string    `json:"location,omitempty" toml:"location,omitempty"`
	URL              string    `json:"url,omitempty" toml:"url,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty" toml:"created_at"`
	UpdatedAt        time.Time `json:"updated_at,omitempty" toml:"updated_at"`
	Active           bool      `json:"active,omitempty" toml:"active"`
	StartsAt         time.Time `json:"starts_at,omitempty" toml:"starts_at"`
	EndsAt           time.Time `json:"ends_at,omitempty" toml:"ends_at"`
	Team             any       `json:"team,omitempty" toml:"team,omitempty"`
}

// JobIn schema
type JobIn struct {
	OrganizationName string `json:"organization_name"`
	Title            string `json:"title"`
	Location         string `json:"location"`
	URL              string `json:"url"`
	Active           bool   `json:"active"`
	Team             any    `json:"team"`
}

// GetJob with given id
func (c *Jobs) GetJob(id int) (out *JobOut, err error) {
	body, err := c.call("GET", fmt.Sprintf("/jobs/%d", id), nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// CreateJob with given payload
func (c *Jobs) CreateJob(in *JobIn) (out *JobOut, err error) {
	body, err := c.call("POST", "/jobs/", in)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}

// UpdateJob with given id and payload
func (c *Jobs) UpdateJob(id int, in *JobIn) (out *JobOut, err error) {
	body, err := c.call("PUT", fmt.Sprintf("/jobs/%d", id), in)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}
