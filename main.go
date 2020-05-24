package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Generated with https://mholt.github.io/json-to-go/
type Releases []struct {
	TagName string `json:"tag_name"`
}

var releaseURL = "https://api.github.com/repos/kubernetes/kubernetes/releases"
var tarballURL = "https://api.github.com/repos/kubernetes/kubernetes/tarball"
var installLocation = "/usr/local/bin/kubectl"

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
	fmt.Print("\n")
	fmt.Printf("Version selected for download is: %v\n", versionWanted)
	fmt.Print("\n")
	downloadFile(installLocation, versionWanted)
	fmt.Println("Downloading Kubernetes version....", versionWanted, "....to", installLocation)
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
	fmt.Printf("Other releases available are: %s\n", data)
	defer resp.Body.Close()
}

func downloadFile(installDirectory string, versionWanted string) {
	resp, err := http.Get(tarballURL + "/" + versionWanted)

	if err != nil {
		fmt.Println("Broken TarBall path - did you select a valid version?", err, http.StatusInternalServerError)
	}

	out, err := os.Create(installDirectory)
	if err != nil {
		fmt.Println("Cannot create file location for kubectl", err)
	}

	_, err = io.Copy(out, resp.Body)

	defer out.Close()
	fmt.Println(_)

}
