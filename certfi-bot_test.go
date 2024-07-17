package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func TestMain(m *testing.M) {
	contents := readConfiguration("config.yml")
	parseConfiguration(contents)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreateToot(t *testing.T) {

	item := gofeed.Item{
		Title:       "This is the title",
		Description: "This is the content",
		Link:        "http://www.suomi24.fi",
	}
	toot := createToot(item)

	if toot.Language != "en" {
		t.Fatalf("Language was not en")
	}
	if toot.Visibility != "public" {
		t.Fatalf("Visibility was not public")
	}
	if toot.Status != fmt.Sprintf("%s\n%s\n\n%s", item.Title, item.Description, item.Link) {
		t.Fatalf("Status was not " + fmt.Sprintf("%s\n%s\n\n%s", item.Title, item.Description, item.Link))
	}
}

func TestIsNotNew(t *testing.T) {
	interval := time.Duration(1) * time.Hour
	pubTime, _ := time.Parse(time.RFC1123, "Wed, 03 Jul 2024 04:10:44 GMT")
	if isNew(pubTime, interval) {
		t.Fatalf("Feed was detected as new while shouldn't have been.")
	}
}

func TestIsNew(t *testing.T) {
	interval := time.Duration(1) * time.Hour
	pubTime := time.Now().Add(time.Duration(-30) * time.Minute)
	if !isNew(pubTime, interval) {
		t.Fatalf("Feed was not detected as new while should have been.")
	}
}

func TestCreateRequest(t *testing.T) {
	expected := toot{
		Status:     "Hello!",
		Visibility: "public",
		Language:   "en",
	}
	request := createRequest(expected)

	authHeader := request.Header.Get("Authorization")
	idempotencyHeader := request.Header.Get("Idempotency-Key")

	if authHeader != "Bearer "+config.Server.AccessToken {
		t.Fatalf("Authorization header was not what is expected")
	}

	if idempotencyHeader == "" {
		t.Fatalf("Idempotency-Key header was not set")
	}

}
