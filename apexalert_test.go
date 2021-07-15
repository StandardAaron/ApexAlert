package main

import "testing"
import "strings"

func TestApiUrlCorrectFormat(t *testing.T) {
	url := getApiUrl()
	if !strings.HasPrefix(url, "https://api.mozambiquehe.re/maprotation?version=2&auth=") {
		t.Fatal("failed")
	}
}
