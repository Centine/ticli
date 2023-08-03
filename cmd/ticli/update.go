package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func checkForUpdates() {
	resp, err := http.Get("https://api.github.com/repos/centine/ticli/releases/latest")
	if err != nil {
		// Handle error
		log.Fatalf("error checking for updates: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("Update notification system not configured")
		return
	} else if resp.StatusCode != http.StatusOK {
		log.Printf("unknown error checking for updates, status code : %d", resp.StatusCode)
	}

	var r Release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatalf("error checking for updates: %s", err)
		return
	}
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(respDump))

	if r.TagName != Version {
		fmt.Printf("A new version is available: %s\n", r.TagName)
	}
}
