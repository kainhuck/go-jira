package v92

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/andygrunwald/go-jira/v2/onpremise"
)

var client *onpremise.Client

func init() {
	var err error

	var gDial = &net.Dialer{
		Timeout: 10 * time.Second,
	}

	c := (&onpremise.BasicAuthTransport{
		Username: "admin",
		Password: "op@ms2022",
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
			DialContext:       gDial.DialContext,
			Dial:              gDial.Dial,
		},
	}).Client()

	client, err = onpremise.NewClient("http://192.168.180.169:8080", c)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func TestGetProjectList(t *testing.T) {
	projectList, _, err := client.Project.GetAll(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, pro := range *projectList {
		fmt.Println(pro.Name, pro.ID)
	}
}

func TestGetProjectIssueTypes(t *testing.T) {
	types, _, err := client.Issue.GetProjectIssueTypes(context.Background(), "10002")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(types)
	for _, m := range types {
		fmt.Println(m.Name, m.ID)
	}
}

func TestGetIssueFields(t *testing.T) {
	fields, _, err := client.Issue.GetProjectIssueFields(context.Background(), "10002", "10100")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fields)
	for _, m := range fields {
		fmt.Println(m.Name, m.FieldID)
	}
}

func TestGetGroups(t *testing.T) {
	groups, _, err := client.Group.GetAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(groups)
	for _, m := range groups {
		fmt.Println(m.Name)
	}
}

// NOTCHANGE
func TestGetGroupMembers(t *testing.T) {
	groups, _, err := client.Group.Get(context.Background(), "jira-administrators", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(groups)
}

// NOTCHANGE
func TestExport(t *testing.T) {
	field := map[string]interface{}{
		"issuetype": map[string]interface{}{
			"id": "10100",
		},
		"project": map[string]interface{}{
			"id": "10002",
		},
		"summary": "test_name_kain_22",
	}

	_, _, err := client.Issue.Create(context.Background(), &onpremise.Issue{
		Fields: &onpremise.IssueFields{
			Unknowns: field,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("OK")
}

// NOTCHANGE
func TestGetSelf(t *testing.T) {
	_, rsp, err := client.User.GetSelf(context.Background())
	if err != nil {
		if _, ok := err.(net.Error); ok {
			log.Fatal("无法和jira服务器建立连接")
		}
		if rsp != nil {
			if rsp.StatusCode == http.StatusForbidden {
				log.Fatalf("jira登录错误，错误码: %v \n"+
					"JIRA错误提示: 账号or密码错误，已经触发JIRA服务的验证码机制，请从web端登出JIRA，"+
					"再登入时填入验证码。操作后才可进行JIRA同步配置操作", rsp.StatusCode)
			}
		}

		log.Fatal("连接jira服务器出错或用户名密码错误")
	}

	if rsp.StatusCode != http.StatusOK {
		log.Fatalf("jira登录错误，错误码: %v", rsp.StatusCode)
	}
	log.Println("SUCCESS")
}
