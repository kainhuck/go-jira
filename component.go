package jira

import (
	"context"
	"net/http"
)

// ComponentService handles components for the Jira instance / API.//
// Jira API docs: https://docs.atlassian.com/software/jira/docs/api/REST/7.10.1/#api/2/component
type ComponentService service

// CreateComponentOptions are passed to the ComponentService.Create function to create a new Jira component
type CreateComponentOptions struct {
	Name         string `json:"name,omitempty" structs:"name,omitempty"`
	Description  string `json:"description,omitempty" structs:"description,omitempty"`
	Lead         *User  `json:"lead,omitempty" structs:"lead,omitempty"`
	LeadUserName string `json:"leadUserName,omitempty" structs:"leadUserName,omitempty"`
	AssigneeType string `json:"assigneeType,omitempty" structs:"assigneeType,omitempty"`
	Assignee     *User  `json:"assignee,omitempty" structs:"assignee,omitempty"`
	Project      string `json:"project,omitempty" structs:"project,omitempty"`
	ProjectID    int    `json:"projectId,omitempty" structs:"projectId,omitempty"`
}

// Create creates a new Jira component based on the given options.
func (s *ComponentService) Create(ctx context.Context, options *CreateComponentOptions) (*ProjectComponent, *Response, error) {
	apiEndpoint := "rest/api/2/component"
	req, err := s.client.NewRequest(ctx, http.MethodPost, apiEndpoint, options)
	if err != nil {
		return nil, nil, err
	}

	component := new(ProjectComponent)
	resp, err := s.client.Do(req, component)

	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}

	return component, resp, nil
}
