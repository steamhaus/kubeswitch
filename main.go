package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

//Generated with https://mholt.github.io/json-to-go/
type Releases []struct {
	TagName string `json:"tag_name"`
}

const (
	stableURL              = "https://storage.googleapis.com/kubernetes-release/release/stable.txt"
	releaseURL             = "https://api.github.com/repos/kubernetes/kubernetes/releases"
	downloadURL            = "https://storage.googleapis.com/kubernetes-release/release/"
	installLocation        = "/usr/local/bin/kubectl"
	binPathLinux           = "/bin/linux/amd64/kubectl"
	binPathMac             = "/bin/darwin/amd64/kubectl"
	GOOS            string = runtime.GOOS
)

func main() {
	fmt.Println("You're using OS: ", GOOS)
	resp, err := http.Get(stableURL)

	if err != nil {
		fmt.Println("Cannot read latest stable version from remote repository", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	result := string(body)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Latest stable release is:" + " " + result + "" + "Do you want to install this version?")
	text, _ := reader.ReadString('\n')

	if strings.TrimRight(text, "\n") == "yes" || strings.TrimRight(text, "\n") == "y" {
		fmt.Println("Downloading Kubernetes version: " + " " + result + " " + "to" + " " + installLocation)
		// There is a bug somewhere appending a new line to the result, causing a nil pointer reference
		downloadFile(installLocation, strings.TrimRight(result, "\n"))

		fmt.Println("version" + " " + result + "has been installed")
		os.Exit(1)
	} else {
		getAllReleases()
		fmt.Println("Which version would you like to install?")
		versionInput, _ := reader.ReadString('\n')
		versionWanted := strings.TrimRight(versionInput, "\n")
		downloadFile(installLocation, versionWanted)
		fmt.Println("Downloading Kubernetes version....", versionWanted, "....to", installLocation)
	}
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
	var binPath string

	if GOOS == "linux" {
		binPath = binPathLinux
	} else if GOOS == "darwin" {
		binPath = binPathMac
	} else {
		os.Exit(0)
	}

	fmt.Println(binPath)

	resp, err := http.Get(downloadURL + versionWanted + binPath)

	fmt.Println(downloadURL + versionWanted + binPath)

	out, err := os.Create("kubectl")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Body)

	n, err := io.Copy(out, resp.Body)
	err = os.Chmod("kubectl", 755)
	if err != nil {
		fmt.Println(err, n)
	}

	x := os.Rename("kubectl", installLocation)
	if x != nil {
		fmt.Println(x)
	}
	defer out.Close()
	defer resp.Body.Close()

}
