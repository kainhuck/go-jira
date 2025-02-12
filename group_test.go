package jira

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestGroupService_GetPage(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/group/member", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, "/rest/api/2/group/member?groupname=default")
		startAt := r.URL.Query().Get("startAt")
		if startAt == "0" {
			fmt.Fprint(w, `{"self":"http://www.example.com/jira/rest/api/2/group/member?includeInactiveUsers=false&maxResults=2&groupname=default&startAt=0","nextPage":"`+testServer.URL+`/rest/api/2/group/member?groupname=default&includeInactiveUsers=false&maxResults=2&startAt=2","maxResults":2,"startAt":0,"total":4,"isLast":false,"values":[{"self":"http://www.example.com/jira/rest/api/2/user?username=michael","name":"michael","key":"michael","emailAddress":"michael@example.com","displayName":"MichaelScofield","active":true,"timeZone":"Australia/Sydney"},{"self":"http://www.example.com/jira/rest/api/2/user?username=alex","name":"alex","key":"alex","emailAddress":"alex@example.com","displayName":"AlexanderMahone","active":true,"timeZone":"Australia/Sydney"}]}`)
		} else if startAt == "2" {
			fmt.Fprint(w, `{"self":"http://www.example.com/jira/rest/api/2/group/member?includeInactiveUsers=false&maxResults=2&groupname=default&startAt=2","maxResults":2,"startAt":2,"total":4,"isLast":true,"values":[{"self":"http://www.example.com/jira/rest/api/2/user?username=michael","name":"michael","key":"michael","emailAddress":"michael@example.com","displayName":"MichaelScofield","active":true,"timeZone":"Australia/Sydney"},{"self":"http://www.example.com/jira/rest/api/2/user?username=alex","name":"alex","key":"alex","emailAddress":"alex@example.com","displayName":"AlexanderMahone","active":true,"timeZone":"Australia/Sydney"}]}`)
		} else {
			t.Errorf("startAt %s", startAt)
		}
	})
	if page, resp, err := testClient.Group.Get(context.Background(), "default", &GroupSearchOptions{
		StartAt:              0,
		MaxResults:           2,
		IncludeInactiveUsers: false,
	}); err != nil {
		t.Errorf("Error given: %s %s", err, testServer.URL)
	} else if page == nil || len(page) != 2 {
		t.Error("Expected members. Group.Members is not 2 or is nil")
	} else {
		if resp.StartAt != 0 {
			t.Errorf("Expect Result StartAt to be 0, but is %d", resp.StartAt)
		}
		if resp.MaxResults != 2 {
			t.Errorf("Expect Result MaxResults to be 2, but is %d", resp.MaxResults)
		}
		if resp.Total != 4 {
			t.Errorf("Expect Result Total to be 4, but is %d", resp.Total)
		}
		if page, resp, err := testClient.Group.Get(context.Background(), "default", &GroupSearchOptions{
			StartAt:              2,
			MaxResults:           2,
			IncludeInactiveUsers: false,
		}); err != nil {
			t.Errorf("Error give: %s %s", err, testServer.URL)
		} else if page == nil || len(page) != 2 {
			t.Error("Expected members. Group.Members is not 2 or is nil")
		} else {
			if resp.StartAt != 2 {
				t.Errorf("Expect Result StartAt to be 2, but is %d", resp.StartAt)
			}
			if resp.MaxResults != 2 {
				t.Errorf("Expect Result MaxResults to be 2, but is %d", resp.MaxResults)
			}
			if resp.Total != 4 {
				t.Errorf("Expect Result Total to be 4, but is %d", resp.Total)
			}
		}
	}
}

func TestGroupService_Add(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/group/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testRequestURL(t, r, "/rest/api/2/group/user?groupname=default")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"name":"default","self":"http://www.example.com/jira/rest/api/2/group?groupname=default","users":{"size":1,"items":[],"max-results":50,"start-index":0,"end-index":0},"expand":"users"}`)
	})

	if group, _, err := testClient.Group.Add(context.Background(), "default", "theodore"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if group == nil {
		t.Error("Expected group. Group is nil")
	}
}

func TestGroupService_Remove(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/group/user", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testRequestURL(t, r, "/rest/api/2/group/user?groupname=default")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"name":"default","self":"http://www.example.com/jira/rest/api/2/group?groupname=default","users":{"size":1,"items":[],"max-results":50,"start-index":0,"end-index":0},"expand":"users"}`)
	})

	if _, err := testClient.Group.Remove(context.Background(), "default", "theodore"); err != nil {
		t.Errorf("Error given: %s", err)
	}
}
