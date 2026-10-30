package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"github.com/victorrbravo/stories-service/constants"
)

func outputToString(output []byte) string {
	return strings.Trim(string(output), "\n ")
}

func getRevCount(binaryPath string) (int, error) {
	cmd := exec.Command(binaryPath, "rev-list", "--count", "HEAD")
	revision, err := cmd.Output()
	if err != nil {
		return -1, err
	}
	text := outputToString(revision)
	value, err := strconv.Atoi(text)
	if err != nil {
		return -1, err
	}
	return value, nil
}

func getCurrentBranch(binaryPath string) (string, error) {
	cmd := exec.Command(binaryPath, "rev-parse", "--abbrev-ref", "HEAD")
	hash, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return outputToString(hash), nil
}

func getHash(binaryPath string) (string, error) {
	cmd := exec.Command(binaryPath, "rev-parse", "--short", "HEAD")
	hash, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return outputToString(hash), nil
}

func main() {
	var workingMode string
	workingMode = strings.Trim(os.Getenv(constants.SystemWorkingMode), " ")
	if len(workingMode) == 0 {
		workingMode = constants.WorkingModeNormal
	}
	parser := template.New("version.go.in")
	binaryPath, err := exec.LookPath("git")
	if err != nil {
		panic(err)
	}
	revision, err := getRevCount(binaryPath)
	if err != nil {
		panic(err)
	}
	hash, err := getHash(binaryPath)
	if err != nil {
		panic(err)
	}
	branch, err := getCurrentBranch(binaryPath)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("../finvar/constants/version.go")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()
	tpl := template.Must(parser.ParseFiles("../generators/version/version.go.in"))
	err = tpl.Execute(file, struct {
		Branch      string
		Revision    int
		Hash        string
		WorkingMode string
	}{
		Branch:      branch,
		Revision:    revision,
		Hash:        hash,
		WorkingMode: workingMode,
	})
	if err != nil {
		panic(err)
	}
}
