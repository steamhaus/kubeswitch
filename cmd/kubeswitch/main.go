package main

import (
  "fmt"
  "io"
  "net/http"
  "os"
  "path/filepath"
	"regexp"
	"strings"
)

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
  fmt.Println(result)
}
