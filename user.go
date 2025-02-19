package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// UserService handles users for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-group-Users
type UserService service

// User represents a Jira user.
type User struct {
	Self            string     `json:"self,omitempty" structs:"self,omitempty"`
	AccountID       string     `json:"accountId,omitempty" structs:"accountId,omitempty"`
	AccountType     string     `json:"accountType,omitempty" structs:"accountType,omitempty"`
	Name            string     `json:"name,omitempty" structs:"name,omitempty"`
	Key             string     `json:"key,omitempty" structs:"key,omitempty"`
	Password        string     `json:"-"`
	EmailAddress    string     `json:"emailAddress,omitempty" structs:"emailAddress,omitempty"`
	AvatarUrls      AvatarUrls `json:"avatarUrls,omitempty" structs:"avatarUrls,omitempty"`
	DisplayName     string     `json:"displayName,omitempty" structs:"displayName,omitempty"`
	Active          bool       `json:"active,omitempty" structs:"active,omitempty"`
	TimeZone        string     `json:"timeZone,omitempty" structs:"timeZone,omitempty"`
	Locale          string     `json:"locale,omitempty" structs:"locale,omitempty"`
	ApplicationKeys []string   `json:"applicationKeys,omitempty" structs:"applicationKeys,omitempty"`
}

// UserGroup represents the group list
type UserGroup struct {
	Self string `json:"self,omitempty" structs:"self,omitempty"`
	Name string `json:"name,omitempty" structs:"name,omitempty"`
}

type userSearchParam struct {
	name  string
	value string
}

type userSearch []userSearchParam

type userSearchF func(userSearch) userSearch

// Get gets user info from Jira using its Account Id
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-user-get
func (s *UserService) Get(ctx context.Context, accountId string) (*User, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user?accountId=%s", accountId)
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return user, resp, nil
}

// GetByAccountID gets user info from Jira
// Searching by another parameter that is not accountId is deprecated,
// but this method is kept for backwards compatibility
// Jira API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-getUser
func (s *UserService) GetByAccountID(ctx context.Context, accountID string) (*User, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user?accountId=%s", accountID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return user, resp, nil
}

// Create creates an user in Jira.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-createUser
func (s *UserService) Create(ctx context.Context, user *User) (*User, *Response, error) {
	apiEndpoint := "/rest/api/2/user"
	req, err := s.client.NewRequest(ctx, http.MethodPost, apiEndpoint, user)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	responseUser := new(User)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		e := fmt.Errorf("could not read the returned data")
		return nil, resp, NewJiraError(resp, e)
	}
	err = json.Unmarshal(data, responseUser)
	if err != nil {
		e := fmt.Errorf("could not unmarshall the data into struct")
		return nil, resp, NewJiraError(resp, e)
	}
	return responseUser, resp, nil
}

// Delete deletes an user from Jira.
// Returns http.StatusNoContent on success.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-user-delete
// Caller must close resp.Body
func (s *UserService) Delete(ctx context.Context, accountId string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user?accountId=%s", accountId)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, NewJiraError(resp, err)
	}
	return resp, nil
}

// GetGroups returns the groups which the user belongs to
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-user-groups-get
func (s *UserService) GetGroups(ctx context.Context, accountId string) (*[]UserGroup, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user/groups?accountId=%s", accountId)
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	userGroups := new([]UserGroup)
	resp, err := s.client.Do(req, userGroups)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return userGroups, resp, nil
}

// GetSelf information about the current logged-in user
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-myself-get
func (s *UserService) GetSelf(ctx context.Context) (*User, *Response, error) {
	const apiEndpoint = "rest/api/2/myself"
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	var user User
	resp, err := s.client.Do(req, &user)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return &user, resp, nil
}

// WithMaxResults sets the max results to return
func WithMaxResults(maxResults int) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "maxResults", value: fmt.Sprintf("%d", maxResults)})
		return s
	}
}

// WithStartAt set the start pager
func WithStartAt(startAt int) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "startAt", value: fmt.Sprintf("%d", startAt)})
		return s
	}
}

// WithActive sets the active users lookup
func WithActive(active bool) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "includeActive", value: fmt.Sprintf("%t", active)})
		return s
	}
}

// WithInactive sets the inactive users lookup
func WithInactive(inactive bool) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "includeInactive", value: fmt.Sprintf("%t", inactive)})
		return s
	}
}

// WithUsername sets the username to search
func WithUsername(username string) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "username", value: username})
		return s
	}
}

// WithAccountId sets the account id to search
func WithAccountId(accountId string) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "accountId", value: accountId})
		return s
	}
}

// WithProperty sets the property (Property keys are specified by path) to search
func WithProperty(property string) userSearchF {
	return func(s userSearch) userSearch {
		s = append(s, userSearchParam{name: "property", value: property})
		return s
	}
}

// Find searches for user info from Jira:
// It can find users by email or display name using the query parameter
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-user-search-get
func (s *UserService) Find(ctx context.Context, property string, tweaks ...userSearchF) ([]User, *Response, error) {
	search := []userSearchParam{
		{
			name:  "query",
			value: property,
		},
	}
	for _, f := range tweaks {
		search = f(search)
	}

	var queryString = ""
	for _, param := range search {
		queryString += param.name + "=" + param.value + "&"
	}

	apiEndpoint := fmt.Sprintf("/rest/api/2/user/search?%s", queryString[:len(queryString)-1])
	req, err := s.client.NewRequest(ctx, http.MethodGet, apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	users := []User{}
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return users, resp, nil
}
