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

//Rune2e executes the script which builds and run the Docker with e2e code
func Rune2e(buildEnv, gitSHA, reportURL string) ([]byte, error) {
	cmd := exec.Command("./script.sh", buildEnv, gitSHA, reportURL)
	fmt.Println(cmd.Args)
	cmd.Dir = currentDir
	stdoutStderr, err := cmd.CombinedOutput()
	return stdoutStderr, err
}
