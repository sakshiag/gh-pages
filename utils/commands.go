package utils

import (
	"fmt"
	"os/exec"
	"path"
	"runtime"
)

var currentDir string

func init() {
	_, filename, _, _ := runtime.Caller(1)
	currentDir = path.Dir(filename)
}

//RunUAT executes the script which builds and run the Docker with UAT code
func RunUAT(buildEnv, gitSHA, reportURL string) ([]byte, error) {
	cmd := exec.Command("./script.sh", buildEnv, gitSHA, reportURL)
	fmt.Println(cmd.Args)
	cmd.Dir = currentDir
	stdoutStderr, err := cmd.CombinedOutput()
	return stdoutStderr, err
}
