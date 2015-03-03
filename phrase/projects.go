package phrase

import (
	"fmt"
)

// ProjectsService provides access to the projects related functions
// in the Phrase API.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/projects/
type ProjectsService struct {
	client *Client
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Get details for the current project.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/projects/
func (s *ProjectsService) Current() (*Project, error) {
	req, err := s.client.NewRequest("GET", "projects/current", nil)
	if err != nil {
		return nil, err
	}

	project := new(Project)
	_, err = s.client.Do(req, project)
	if err != nil {
		return nil, err
	}

	return project, err
}

func (p Project) String() string {
	return fmt.Sprintf("Project ID: %d Name: %s",
		p.ID, p.Name)
}
