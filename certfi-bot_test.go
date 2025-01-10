package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pbloigu/gonfig"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gonfig.Get("config.yml", &config)
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

	assert.Equal(t, "en", toot.Language, "Language wrong.")
	assert.Equal(t, "public", toot.Visibility, "Visibility wrong.")
	assert.Equal(t, fmt.Sprintf("%s\n%s\n\n%s", item.Title, item.Description, item.Link), toot.Status, "Status wrong.")

}

func TestIsNotNew(t *testing.T) {
	interval := time.Duration(1) * time.Hour
	pubTime, _ := time.Parse(time.RFC1123, "Wed, 03 Jul 2024 04:10:44 GMT")
	assert.False(t, isNew(pubTime, interval), "Feed was detected as new while shouldn't have been.")
}

func TestIsNew(t *testing.T) {
	interval := time.Duration(1) * time.Hour
	pubTime := time.Now().Add(time.Duration(-30) * time.Minute)
	assert.True(t, isNew(pubTime, interval), "Feed was not detected as new while should have been.")
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

	assert.Equal(t, "Bearer "+config.Server.AccessToken, authHeader, "Authorization header was not what is expected")
	assert.NotNil(t, idempotencyHeader, "Idempotency-Key header was not set")
}
