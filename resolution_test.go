package jira

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestResolutionService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/rest/api/2/resolution"

	raw, err := os.ReadFile("../testing/mock-data/all_resolutions.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	resolution, _, err := testClient.Resolution.GetList(context.Background())
	if resolution == nil {
		t.Error("Expected resolution list. Resolution list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
