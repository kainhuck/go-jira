package jira

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// GroupService handles Groups for the Jira instance / API.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/server/#api/2/group
type GroupService service

// groupMembersResult is only a small wrapper around the Group* methods
// to be able to parse the results
type groupMembersResult struct {
	StartAt    int           `json:"startAt"`
	MaxResults int           `json:"maxResults"`
	Total      int           `json:"total"`
	Members    []GroupMember `json:"values"`
}

// Group represents a Jira group
type Group struct {
	ID                   string          `json:"id"`
	Title                string          `json:"title"`
	Type                 string          `json:"type"`
	Properties           groupProperties `json:"properties"`
	AdditionalProperties bool            `json:"additionalProperties"`
}

type groupProperties struct {
	Name groupPropertiesName `json:"name"`
}

type groupPropertiesName struct {
	Type string `json:"type"`
}

// GroupMember reflects a single member of a group
type GroupMember struct {
	Self         string `json:"self,omitempty"`
	Name         string `json:"name,omitempty"`
	Key          string `json:"key,omitempty"`
	AccountID    string `json:"accountId,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	Active       bool   `json:"active,omitempty"`
	TimeZone     string `json:"timeZone,omitempty"`
	AccountType  string `json:"accountType,omitempty"`
}

type GroupMembersUsers struct {
	Size         int           `json:"size"`
	MaxResults   int           `json:"max-results"`
	StartIndex   int           `json:"start-index"`
	EndIndex     int           `json:"end-index"`
	GroupMembers []GroupMember `json:"items"`
}

type GetGroupMembersResult struct {
	Name   string            `json:"name"`
	Self   string            `json:"self"`
	Expand string            `json:"expand"`
	Users  GroupMembersUsers `json:"users"`
}

type GetGroupMembersResultV9 struct {
	Self       string        `json:"self"`
	NextPage   string        `json:"nextPage"`
	MaxResults int           `json:"maxResults"`
	StartAt    int           `json:"startAt"`
	Total      int           `json:"total"`
	IsLast     bool          `json:"isLast"`
	Users      []GroupMember `json:"values"`
}

// GroupMeta JIRA v92
type GroupMeta struct {
	Name   string `json:"name,omitempty"`
	HTML   string `json:"html,omitempty"`
	Labels []struct {
		Text  string `json:"text,omitempty"`
		Title string `json:"title,omitempty"`
		Type  string `json:"type,omitempty"`
	} `json:"labels,omitempty"`
}

// GroupsResult JIRA v9.2
type GroupsResult struct {
	Header string      `json:"header,omitempty"`
	Total  int         `json:"total"`
	Groups []GroupMeta `json:"groups"`
}

// GroupSearchOptions specifies the optional parameters for the Get Group methods
type GroupSearchOptions struct {
	StartAt              int
	MaxResults           int
	IncludeInactiveUsers bool
}

// Get returns a paginated list of members of the specified group and its subgroups.
// Users in the page are ordered by user names.
// User of this resource is required to have sysadmin or admin permissions.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/server/#api/2/group-getUsersFromGroup
//
// WARNING: This API only returns the first page of group members
func (s *GroupService) Get(ctx context.Context, name string, options *GroupSearchOptions) ([]GroupMember, *Response, error) {
	var apiEndpoint string
	if options == nil {
		apiEndpoint = fmt.Sprintf("/rest/api/2/group/member?groupname=%s", url.QueryEscape(name))
	} else {
		// TODO use addOptions
		apiEndpoint = fmt.Sprintf(
			"/rest/api/2/group/member?groupname=%s&startAt=%d&maxResults=%d&includeInactiveUsers=%t",
			url.QueryEscape(name),
			options.StartAt,
			options.MaxResults,
			options.IncludeInactiveUsers,
		)
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	group := new(groupMembersResult)
	resp, err := s.client.Do(req, group)
	if err != nil {
		return nil, resp, err
	}
	return group.Members, resp, nil
}

func (s *GroupService) GetAll(ctx context.Context) ([]GroupMeta, *Response, error) {
	apiEndpoint := "rest/api/2/groups/picker?maxResults=1000000"
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	meta := new(GroupsResult)
	resp, err := s.client.Do(req, meta)

	if err != nil {
		return nil, resp, err
	}

	return meta.Groups, resp, nil
}

func (s *GroupService) GetGroupMembersV9(ctx context.Context, name string) ([]GroupMember, *Response, error) {
	var (
		index = 0
		all   = make([]GroupMember, 0)
		req   *http.Request
		resp  *Response
		err   error
	)

	for {
		apiEndpoint := fmt.Sprintf("/rest/api/2/group/member?groupname=%s&startAt=%d", name, index)
		req, err = s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
		if err != nil {
			return nil, nil, err
		}

		group := new(GetGroupMembersResultV9)
		resp, err = s.client.Do(req, group)
		if err != nil {
			return nil, resp, err
		}

		all = append(all, group.Users...)
		if group.IsLast {
			break
		}
		index += 50
	}

	return all, resp, nil
}

func (s *GroupService) GetGroupMembers(ctx context.Context, name string) ([]GroupMember, *Response, error) {
	users, resp, err := s.GetGroupMembersV9(ctx, name)
	if err == nil {
		return users, resp, nil
	}
	apiEndpoint := fmt.Sprintf("rest/api/2/group?groupname=%s&expand=users", url.QueryEscape(name))
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	group := new(GetGroupMembersResult)
	resp, err = s.client.Do(req, group)
	if err != nil {
		return nil, resp, err
	}

	return group.Users.GroupMembers, resp, nil
}

// Add adds user to group
//
// Jira API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/group-addUserToGroup
func (s *GroupService) Add(ctx context.Context, groupname string, username string) (*Group, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/group/user?groupname=%s", groupname)
	var user struct {
		Name string `json:"name"`
	}
	user.Name = username
	req, err := s.client.NewRequest(ctx, http.MethodPost, apiEndpoint, &user)
	if err != nil {
		return nil, nil, err
	}

	responseGroup := new(Group)
	resp, err := s.client.Do(req, responseGroup)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return responseGroup, resp, nil
}

// Remove removes user from group
//
// Jira API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/group-removeUserFromGroup
// Caller must close resp.Body
func (s *GroupService) Remove(ctx context.Context, groupname string, username string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/group/user?groupname=%s&username=%s", groupname, username)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return resp, jerr
	}

	return resp, nil
}
