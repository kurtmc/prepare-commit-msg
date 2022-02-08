package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	var dat []byte
	var err error
	dat, err = ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	branchName, err := getCurrentBranch()
	if err != nil {
		return
	}
	err = ioutil.WriteFile(os.Args[1], []byte(prepareMessage(branchName, string(dat))), 0644)
	if err != nil {
		panic(err)
	}
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	branch := strings.TrimSpace(fmt.Sprintf("%s", out))
	return branch, nil
}

func prepareMessage(branch, message string) string {
	branchRegexp := regexp.MustCompile("[A-Z0-9]+-[0-9]+")
	messageRegexp := regexp.MustCompile("^[A-Z0-9]+-[0-9]+: ")

	if !branchRegexp.MatchString(branch) {
		return message
	}

	matches := branchRegexp.FindAllStringSubmatch(branch, 1)
	ticket := matches[0][0]

	messagePrefix := ticket + ": "

	if messageRegexp.MatchString(message) {
		messageBranch := branchRegexp.FindString(message)
		if messageBranch == branch {
			return message
		}
		return messageRegexp.ReplaceAllString(message, messagePrefix)
	}

	return messagePrefix + message
}
