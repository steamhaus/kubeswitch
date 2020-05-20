package main

import (
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

type kubeVersion struct {
	version string
	url     string
}

func main() {

	resp, err := http.Get("https://storage.googleapis.com/kubernetes-release/release/stable.txt")

	if err != nil {
		fmt.Println("Cannot read from remote repository")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	result := string(body)
	fmt.Println("Latest stable release is:" + " " + result)
}
