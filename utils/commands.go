package utils

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

//It will clone the git repo which contains the configuration file.
func cloneRepo(msg Message) ([]byte, string, error) {
	gitURL := msg.GitURL
	cmd := exec.Command("git", "clone", gitURL)
	fmt.Println(cmd.Args)
	cmd.Dir = currentDir
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", err
	}
	urlPath, err := url.Parse(msg.GitURL)
	if err != nil {
		return nil, "", err
	}
	baseName := filepath.Base(urlPath.Path)
	extName := filepath.Ext(urlPath.Path)
	p := baseName[:len(baseName)-len(extName)]
	path := currentDir + "/" + p + "/terraform.tfvars"
	createFile(msg, path)
	return stdoutStderr, p, err
}

//It will create a vars file
func createFile(msg Message, path string) {
	// detect if file exists

	_, err := os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
	}

	writeFile(path, msg)
}

func writeFile(path string, msg Message) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	variables := msg.VariableStore

	for _, v := range *variables {
		_, err = file.WriteString(v.Name + " = \"" + v.Value + "\" \n")
	}

	// save changes
	err = file.Sync()
	if err != nil {
		return
	}
}
