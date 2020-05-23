package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	/*
		"os"
		"io"
		"path/filepath"
		"regexp"
		"string"
	*/)

// Generated with https://mholt.github.io/json-to-go/
type Releases []struct {

	TagName string `json:"tag_name"`
	// TargetCommitish string `json:"target_commitish"`
	// Name string `json:"name"`
	// Draft           bool   `json:"draft"`
	// Author          struct {
	// 	Login             string `json:"login"`
	// 	ID                int    `json:"id"`
	// 	NodeID            string `json:"node_id"`
	// 	AvatarURL         string `json:"avatar_url"`
	// 	GravatarID        string `json:"gravatar_id"`
	// 	URL               string `json:"url"`
	// 	HTMLURL           string `json:"html_url"`
	// 	FollowersURL      string `json:"followers_url"`
	// 	FollowingURL      string `json:"following_url"`
	// 	GistsURL          string `json:"gists_url"`
	// 	StarredURL        string `json:"starred_url"`
	// 	SubscriptionsURL  string `json:"subscriptions_url"`
	// 	OrganizationsURL  string `json:"organizations_url"`
	// 	ReposURL          string `json:"repos_url"`
	// 	EventsURL         string `json:"events_url"`
	// 	ReceivedEventsURL string `json:"received_events_url"`
	// 	Type              string `json:"type"`
	// 	SiteAdmin         bool   `json:"site_admin"`
	// } `json:"author"`
	// Prerelease  bool      `json:"prerelease"`
	// CreatedAt   time.Time `json:"created_at"`
	// PublishedAt time.Time `json:"published_at"`
	// Assets      []struct {
	// 	URL      string `json:"url"`
	// 	ID       int    `json:"id"`
	// 	// NodeID   string `json:"node_id"`
	// 	Name     string `json:"name"`
	// Label    string `json:"label"`
	// Uploader struct {
	// 	Login             string `json:"login"`
	// 	ID                int    `json:"id"`
	// 	NodeID            string `json:"node_id"`
	// 	AvatarURL         string `json:"avatar_url"`
	// 	GravatarID        string `json:"gravatar_id"`
	// 	URL               string `json:"url"`
	// 	HTMLURL           string `json:"html_url"`
	// 	FollowersURL      string `json:"followers_url"`
	// 	FollowingURL      string `json:"following_url"`
	// 	GistsURL          string `json:"gists_url"`
	// 	StarredURL        string `json:"starred_url"`
	// 	SubscriptionsURL  string `json:"subscriptions_url"`
	// 	OrganizationsURL  string `json:"organizations_url"`
	// 	ReposURL          string `json:"repos_url"`
	// 	EventsURL         string `json:"events_url"`
	// 	ReceivedEventsURL string `json:"received_events_url"`
	// 	Type              string `json:"type"`
	// 	SiteAdmin         bool   `json:"site_admin"`
	// } `json:"uploader"`
	// 	ContentType        string    `json:"content_type"`
	// 	State              string    `json:"state"`
	// 	Size               int       `json:"size"`
	// 	DownloadCount      int       `json:"download_count"`
	// 	CreatedAt          time.Time `json:"created_at"`
	// 	UpdatedAt          time.Time `json:"updated_at"`
	// 	//BrowserDownloadURL string    `json:"browser_download_url"`
	// } `json:"assets"`
	//	TarballURL string `json:"tarball_url"`
	// ZipballURL string `json:"zipball_url"`
	// //Body       string `json:"body"`
}

func main() {

	resp, err := http.Get("https://storage.googleapis.com/kubernetes-release/release/stable.txt")

	if err != nil {
		fmt.Println("Cannot read from remote repository", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	result := string(body)
	fmt.Println("Latest stable release is:" + " " + result)
	getAllReleases()
}

func getAllReleases() {
	releaseURL := "https://api.github.com/repos/kubernetes/kubernetes/releases"

	resp, err := http.Get(releaseURL)

	if err != nil {
		fmt.Println("Cannot see all latest releaes", err, http.StatusInternalServerError)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Read body: %v", err)
	}
	var data Releases
	json.Unmarshal(body, &data)

	//TODO: Work out how we can format this list better wit a new line after each result
	fmt.Printf("Other releases available are: %v\n", data)
}

