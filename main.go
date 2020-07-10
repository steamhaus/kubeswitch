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

	"github.com/akamensky/argparse"
)

//Releases is enerated with https://mholt.github.io/json-to-go/
type Releases []struct {
	TagName string `json:"name"`
}

const (
	stableURL       = "https://storage.googleapis.com/kubernetes-release/release/stable.txt"
	releaseURL      = "https://api.github.com/repos/kubernetes/kubernetes/releases"
	downloadURL     = "https://storage.googleapis.com/kubernetes-release/release/"
	installLocation = "/usr/local/bin/kubectl"
	binPathLinux    = "/bin/linux/amd64/kubectl"
	binPathMac      = "/bin/darwin/amd64/kubectl"

	//GOOS is used to detect the OS used by the host
	GOOS = runtime.GOOS
)

var binPath string
var versionToInstall string

func checkOS() {

	if GOOS == "linux" {
		binPath = binPathLinux
	} else if GOOS == "darwin" {
		binPath = binPathMac
	} else {
		os.Exit(0)
	}
}

func main() {
	checkOS()
	parser := argparse.NewParser("kubeswitch", "easily swap kubectl versions")

	versionFlag := parser.String("v", "version", &argparse.Options{Required: false, Help: "specifiy a version to download"})
	versionToInstall = *versionFlag

	// Maybe in the future there can be an acceptance check, but if we want this as part of an automated sequence it make senses to assume user input is always a correct version.
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}

	if *versionFlag != "" {
		fmt.Println("You've selected kubectl version: ", *versionFlag, "to install")
		fmt.Println("Downloading kubectl version: ", *versionFlag, "to ", installLocation)
		downloadFile(installLocation, *versionFlag)
		fmt.Println(*versionFlag)
	} else if *versionFlag == "" {
		getStable()
	}

}

func getStable() {
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
	fmt.Println("Latest stable release is:" + " " + result + "" + "Do you want to install this version? [yes/no]")
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
	fmt.Printf("Other releases available are: %v\n", data)
	defer resp.Body.Close()
}

func downloadFile(installDirectory string, versionWanted string) {

	resp, err := http.Get(downloadURL + versionWanted + binPath)

	out, err := os.Create("kubectl")

	if err != nil {
		fmt.Println(err)
	}

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
