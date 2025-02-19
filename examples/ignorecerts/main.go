package main

import (
	"context"
	"crypto/tls"
	"fmt"
	jira "github.com/kainhuck/go-jira"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	jiraClient, _ := jira.NewClient("https://issues.apache.org/jira/", client)
	issue, _, _ := jiraClient.Issue.Get(context.Background(), "MESOS-3325", nil)

	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)

}
