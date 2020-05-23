package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Generated with https://mholt.github.io/json-to-go/
type Releases []struct {
	TagName string `json:"tag_name"`
}

type DownloadURL []struct {
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
}

var releaseURL = "https://api.github.com/repos/kubernetes/kubernetes/releases"

func main() {
	versionWanted := os.Args[1]
	resp, err := http.Get("https://storage.googleapis.com/kubernetes-release/release/stable.txt")

	if err != nil {
		fmt.Println("Cannot read from remote repository", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	result := string(body)
	fmt.Println("Latest stable release is:" + " " + result)
	getAllReleases()
	fmt.Printf("Version selected for download is: %v\n", versionWanted)
	fmt.Println("Downloading Kubernetes version....", versionWanted, "....to /usr/bin/kubectl")
}

func getAllReleases() {
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

func getDownloadLocations(installDirectory string) {
	resp, err := http.Get(releaseURL)

}
