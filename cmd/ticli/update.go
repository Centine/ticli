package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func checkForUpdates() {
	resp, err := http.Get("https://api.github.com/repos/centine/ticli/releases/latest")
	if err != nil {
		// Handle error
		return
	}
	defer resp.Body.Close()

	var r Release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		// Handle error
		return
	}

	if r.TagName != version {
		fmt.Printf("A new version is available: %s\n", r.TagName)
	}
}
