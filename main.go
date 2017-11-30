package main

import (
	"gopkg.in/libgit2/git2go.v26"
	"io/ioutil"
	"os"
	"regexp"
)

func main() {
	var dat []byte
	var err error
	dat, err = ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var repo *git.Repository
	repo, err = git.OpenRepository(pwd)
	var ref *git.Reference
	ref, err = repo.Head()
	if err != nil {
		return
	}
	var branchName string
	branchName, err = ref.Branch().Name()
	if err != nil {
		return
	}
	err = ioutil.WriteFile(os.Args[1], []byte(PrepareMessage(branchName, string(dat))), 0644)
	if err != nil {
		panic(err)
	}
}

func PrepareMessage(branch, message string) string {
	var branchRegexp, messageRegexp *regexp.Regexp
	var err error
	branchRegexp, err = regexp.Compile("[A-Z]+-[0-9]+")
	if err != nil {
		panic(err)
	}
	messageRegexp, err = regexp.Compile("^[A-Z]+-[0-9]+: ")
	if err != nil {
		panic(err)
	}

	if !branchRegexp.MatchString(branch) {
		return message
	}

	if messageRegexp.MatchString(message) {
		messageBranch := branchRegexp.FindString(message)
		if messageBranch == branch {
			return message
		} else {
			return messageRegexp.ReplaceAllString(message, branch+": ")
		}
	}

	return branch + ": " + message
}
