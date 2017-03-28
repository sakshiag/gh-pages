package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var httpClient *http.Client
var githubToken string

//GitRef contains the git ref response
type GitRef struct {
	Object struct {
		Sha string `json:"sha"`
	} `json:"object"`
}

const gitAPI = "https://github.ibm.com/api/v3/repos/blueprint/bluemix-terraform-provider-dev/git/refs/heads/master"
const defaultReportURL = "http://9.47.83.184:8080"

func init() {
	httpClient = &http.Client{CheckRedirect: nil}
	githubToken = os.Getenv("GITHUB_API_TOKEN")
	if githubToken == "" {
		panic("GITHUB_API_TOKEN is empty")
	}
}

func headSHA() (string, error) {
	req, _ := http.NewRequest("GET", gitAPI, nil)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", githubToken))
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Couldn't read the git sha. Error is %v", err)
		return "", err
	}
	decoder := json.NewDecoder(res.Body)
	var ref GitRef
	err = decoder.Decode(&ref)
	if err != nil {
		log.Printf("Couldn't decode the git sha. Error is %v", err)
		return "", err
	}
	return ref.Object.Sha, nil
}

//UATHandler handles request to kickoff UAT
func UATHandler(w http.ResponseWriter, r *http.Request) {
	buildEnv := r.Header.Get("BUILD_ENV")
	gitSHA := r.Header.Get("GIT_SHA")
	reportURL := r.Header.Get("REPORT_URL")

	if reportURL == "" {
		reportURL = defaultReportURL
	}
	if buildEnv == "" {
		w.WriteHeader(400)
		w.Write([]byte("Missing BUILD_ENV"))
		return
	}

	if gitSHA == "" {
		log.Println("GIT_SHA not present in the request Header. Will fetch the latest commit")
	}

	if gitSHA == "" {
		sha, err := headSHA()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Couldn't get the git sha %v", err)))
			return
		}
		gitSHA = sha
		log.Println("Will run UAT against", gitSHA)
	}

	go func() {
		output, _ := RunUAT(buildEnv, gitSHA, reportURL)
		fmt.Printf("%s\n", output)
	}()

	io.WriteString(w, "Request to start the UAT submitted succefully")

}
