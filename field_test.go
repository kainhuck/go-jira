package jira

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestFieldService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/rest/api/2/field"

	raw, err := os.ReadFile("../testing/mock-data/all_fields.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	fields, _, err := testClient.Field.GetList(context.Background())
	if fields == nil {
		t.Error("Expected field list. Field list is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
