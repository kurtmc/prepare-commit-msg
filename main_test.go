package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestPrepareMessage(t *testing.T) {
	for _, c := range []struct {
		branch, message, want string
	}{
		{"ARC-123", "My commit", "ARC-123: My commit"},
		{"master", "My commit", "My commit"},
		{"ARC-123", "ARC-123: My commit", "ARC-123: My commit"},
		{"testing", "My commit", "My commit"},
		{"ARC-456", "ARC-123: My commit", "ARC-456: My commit"},
	} {
		got := prepareMessage(c.branch, c.message)
		if got != c.want {
			t.Errorf("PrepareMessage(%q, %q) == %q, want %q", c.branch, c.message, got, c.want)
		}
	}
}

func TestGetCurrentBranch(t *testing.T) {
	dir, err := ioutil.TempDir("", "prepare-commit-msg.test")
	if err != nil {
		t.Fatal("Could not create temp directory", err)
	}

	defer os.RemoveAll(dir) // clean up

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal("Could not change directory to ", dir)
	}

	cmd := exec.Command("bash", "-c", `git init && git config user.email a@a.com && echo "" > a.txt && git add . && git commit -m abc && git checkout -b PLA-123`)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Sprintf("Could not change create git repo: %s", out), err)
	}

	branch, err := getCurrentBranch()
	if err != nil {
		t.Fatal("Could not get current branch: ", err)
	}

	if branch != "PLA-123" {
		t.Fatal("Expected branch to be 'PLA-123' but was '" + branch + "'")
	}
}
