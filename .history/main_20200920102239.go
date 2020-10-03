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
	"os/exec"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/akamensky/argparse"
)

//Releases for Kubernetes follow the name tag
type Releases []struct {
	TagName string `json:"name"`
}

//HelmReleases for Helm follow the tag_name tag
type HelmReleases []struct {
	TagName string `json:"tag_name"`
}

const (
	stableURL           = "https://storage.googleapis.com/kubernetes-release/release/stable.txt"
	releaseURLHelm      = "https://api.github.com/repos/helm/helm/releases?per_page=10"
	releaseURL          = "https://api.github.com/repos/kubernetes/kubernetes/releases?per_page=50"
	downloadURLKube     = "https://storage.googleapis.com/kubernetes-release/release/"
	downloadURLHelm     = "https://get.helm.sh/helm-"
	installLocation     = "/usr/local/bin/kubectl"
	installLocationHelm = "/usr/local/bin/helm"
	binPathLinux        = "/bin/linux/amd64/kubectl"
	binPathMac          = "/bin/darwin/amd64/kubectl"
	helmZipLinux        = "linux-amd64.tar.gz"
	helmZipMac          = "darwin-amd64.tar.gz"

	//GOOS is used to detect the OS used by the host
	GOOS = runtime.GOOS
)

var binPath string
var versionToInstall string
var zipPath string
var downloadUR

func checkOS() {

	if GOOS == "linux" {
		binPath = binPathLinux
		zipPath = helmZipLinux
	} else if GOOS == "darwin" {
		binPath = binPathMac
		zipPath = helmZipMac
	} else {
		os.Exit(1)
	}
}

func main() {
	checkOS()
	parser := argparse.NewParser("kubeswitch", "easily swap kubectl versions")

	versionFlag := parser.String("v", "version", &argparse.Options{Required: false, Help: "specifiy a version to download"})

	awsFlag := parser.Flag("a", "aws", &argparse.Options{Required: false, Help: "Check your AWS EKS/Kops version"})

	helmFlag := parser.Flag("t", "helm", &argparse.Options{Required: false, Help: "Get helm version"})

	// Maybe in the future there can be an acceptance check, but if we want this as part of an automated sequence it make senses to assume user input is always a correct version.
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}

	if *awsFlag {
		checkAWSAuth()
	}

	if *helmFlag {
		getHelm()
		os.Exit(0)
	}

	if *versionFlag != "" {
		fmt.Println("You've selected kubectl version: ", *versionFlag, "to install")
		fmt.Println("Downloading kubectl version: ", *versionFlag, "to ", installLocation)
		downloadFile(installLocation, *versionFlag, "kubectl")
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
		fmt.Println("Downloading Kubernetes version: " + " " + result + "to" + " " + installLocation)
		// There is a bug somewhere appending a new line to the result, causing a nil pointer reference
		downloadFile(installLocation, strings.TrimRight(result, "\n"), "kubectl")

		fmt.Println("version" + " " + result + "has been installed")
		os.Exit(0)
	} else {
		getAllReleases()
		fmt.Println("Which version would you like to install?")
		versionInput, _ := reader.ReadString('\n')
		versionWanted := strings.TrimRight(versionInput, "\n")
		downloadFile(installLocation, versionWanted, "kubectl")
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

func downloadFile(installDirectory string, versionWanted string, app string) {

	if app == "helm" {
		downloadURL = downloadURLHelm + versionWanted + zipPath
	}

	resp, err := http.Get(downloadURL + versionWanted + binPath)

	out, err := os.Create(app)

	if err != nil {
		fmt.Println(err)
	}

	n, err := io.Copy(out, resp.Body)
	err = os.Chmod(app, 755)
	if err != nil {
		fmt.Println(err, n)
	}

	x := os.Rename(app, installLocation)
	if x != nil {
		fmt.Println(x)
	}
	defer out.Close()
	defer resp.Body.Close()

}

func checkAWSAuth() {
	creds := credentials.NewEnvCredentials()

	_, err := creds.Get()
	if err != nil {
		fmt.Println("AWS Credentials not found or set. Skipping...")
	} else {
		out, _ := exec.Command("kubectl", "version", "--short").Output()
		verParse := string(out)
		ver := verParse[40:47]

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Your EKS cluster version is: ", ver, "- do you want to match your client version? [yes/no]")
		text, _ := reader.ReadString('\n')
		if strings.TrimRight(text, "\n") == "yes" {
			downloadFile(installLocation, ver, "kubectl")
			fmt.Println("Client and Server matched.")
			os.Exit(0)
		}
	}

}

func getHelm() {
	reader := bufio.NewReader(os.Stdin)
	resp, err := http.Get(releaseURLHelm)

	if err != nil {
		fmt.Println("Cannot see all latest releaes", err, http.StatusInternalServerError)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Read body: %v", err)
	}
	var data HelmReleases
	json.Unmarshal(body, &data)

	fmt.Printf("Helm releases available are: %v\n", data)
	defer resp.Body.Close()

	fmt.Println("Which version of Helm would you like to install?")
	versionInput, _ := reader.ReadString('\n')
	versionWanted := strings.TrimRight(versionInput, "\n")
	downloadFile(installLocation, versionWanted, "helm")
	fmt.Println("Downloading Helm version....", versionWanted, "....to", installLocationHelm)
}
