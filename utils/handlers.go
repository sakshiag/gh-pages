package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
)

var httpClient *http.Client
var githubToken string
var githubIBMToken string
var planTimeOut = 20 * time.Minute

// Message -
type Message struct {
	GitURL        string            `json:"git_url,required" description:"The git url of your configuraltion"`
	VariableStore *VariablesRequest `json:"variablestore,omitempty" description:"The environments' variable store"`
	LOGLEVEL      string            `json:"log_level,omitempty" description:"The log level defing by user."`
}

// ConfigResponse -
type ConfigResponse struct {
	ID string `json:"id,required" description:"ID of the git operation."`
}

// ActionResponse -
type ActionResponse struct {
	ConfigName string `json:"id,required" description:"Name of the configuration"`
	Output     string `json:"output,required" description:"Output logs of terraform command"`
	Error      string `json:"error,required" description:"Error logs for terraform command"`
	Action     string `json:"action,required" description:"Action Name"`
}

// VariablesRequest -
type VariablesRequest []EnvironmentVariableRequest

// EnvironmentVariableRequest -
type EnvironmentVariableRequest struct {
	Name  string `json:"name,required" binding:"required" description:"The variable's name"`
	Value string `json:"value,required" binding:"required" description:"The variable's value"`
}

var currentDir = "/tmp"

var logDir = "/tmp/log/"

var stateDir = "/tmp/state"

func init() {

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, os.ModePerm)
	}
	if _, err := os.Stat(stateDir); os.IsNotExist(err) {
		os.MkdirAll(stateDir, os.ModePerm)
	}
}

//ConfHandler handles request to kickoff git clone of the repo.
func ConfHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method.", 405)
	}

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg Message
	var response ConfigResponse
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println(msg.GitURL)
	if msg.GitURL == "" {
		w.WriteHeader(400)
		w.Write([]byte("EMPTY GIT URL"))
		return
	}

	if msg.LOGLEVEL != "" {
		os.Setenv("TF_LOG", msg.LOGLEVEL)
	}

	log.Println("Will clone git repo")

	_, id, err := cloneRepo(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println("\n", id)

	response.ID = id
	log.Println(response)

	output, err := json.Marshal(response)
	if err != nil {
		return
	}

	confDir := path.Join(currentDir, id)

	b = make([]byte, 10)
	rand.Read(b)
	randomID := fmt.Sprintf("%x", b)

	err = TerraformInit(confDir, id, &planTimeOut, randomID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

//PlanHandler handles request to run terraform plan.
func PlanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method.", 405)
		return
	}
	vars := mux.Vars(r)
	repoName := vars["repo_name"]

	var response ConfigResponse

	log.Println("Url Param 'repo name' is: " + repoName)
	confDir := path.Join(currentDir, repoName)

	b := make([]byte, 10)
	rand.Read(b)
	randomID := fmt.Sprintf("%x", b)
	go func() {
		err := TerraformPlan(confDir, repoName, &planTimeOut, randomID)
		if err != nil {
			return
		}
	}()
	w.WriteHeader(202)
	response.ID = randomID
	output, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)

}

//ApplyHandler handles request to run terraform plan.
func ApplyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method.", 405)
		return
	}
	var response ConfigResponse
	vars := mux.Vars(r)
	repoName := vars["repo_name"]

	log.Println("Url Param 'repo name' is: " + repoName)
	confDir := path.Join(currentDir, repoName)

	b := make([]byte, 10)
	rand.Read(b)
	randomID := fmt.Sprintf("%x", b)
	go func() {
		err := TerraformApply(confDir, stateDir, repoName, &planTimeOut, randomID)
		if err != nil {
			return
		}
	}()
	w.WriteHeader(202)
	response.ID = randomID
	output, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)

}

//DestroyHandler handles request to run terraform plan.
func DestroyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method.", 405)
		return
	}
	var response ConfigResponse
	vars := mux.Vars(r)
	repoName := vars["repo_name"]

	log.Println("Url Param 'repo name' is: " + repoName)
	confDir := path.Join(currentDir, repoName)

	b := make([]byte, 10)
	rand.Read(b)
	randomID := fmt.Sprintf("%x", b)
	go func() {
		err := TerraformDestroy(confDir, stateDir, repoName, &planTimeOut, randomID)
		if err != nil {
			return
		}
	}()
	w.WriteHeader(202)
	response.ID = randomID
	output, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)

}

//LogHandler handles request to run terraform plan.
func LogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method.", 405)
		return
	}
	var response ActionResponse

	vars := mux.Vars(r)
	repoName := vars["repo_name"]
	action := vars["action"]
	logID := vars["logID"]

	log.Println("Url Param 'repo name' is: " + repoName)
	log.Println("Url Param 'action' is: " + action)
	log.Println("Url Param 'log id' is: " + logID)

	outFile, errFile, err := readLogFile(logID)

	response.ConfigName = repoName
	response.Output = outFile
	response.Error = errFile
	response.Action = action

	output, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)

}
