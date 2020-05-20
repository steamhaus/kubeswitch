package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	/*
		"os"
		"io"
		"path/filepath"
		"regexp"
		"string"
	*/)

type kubeVersion struct {
	version string
	url     string
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
	var kubeVersion kubeVersion
	releaseURL := "https://github.com/kubernetes/kubernetes/tags"

	resp, err := http.Get(releaseURL)

	if err != nil {
		fmt.Println("Cannot see all latest releaes", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	releases := string(body)
	result := strings.Split(releases, "\n")
	for i := range result {
		r, _ := regexp.Compile(`\/(\d+)(\.)(\d+)(\.)(\d+)\/`)

		if r.MatchString(result[i]) {

			str := r.FindString(result[i])
			trimstr := strings.Trim(str, "tag/")

			kubeVersion.version = trimstr
			kubeVersion.url = releaseURL + trimstr
		}

		fmt.Println("All of the latest releases are:"+" ", kubeVersion)
	}
}
