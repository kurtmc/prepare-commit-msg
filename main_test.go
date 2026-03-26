package main

import (
	"fmt"
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
		{"feature/thing/ARC-456", "My commit", "ARC-456: My commit"},
		{"feature/thing/ARC-456", "ARC-456: My commit", "ARC-456: My commit"},
		{"feature/ARC-789/thing", "My commit", "ARC-789: My commit"},
	} {
		got := prepareMessage(c.branch, c.message)
		if got != c.want {
			t.Errorf("PrepareMessage(%q, %q) == %q, want %q", c.branch, c.message, got, c.want)
		}
	}
}

func TestGenerateCommitMessageWithClaude(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	diff := `diff --git a/main.go b/main.go
index abc1234..def5678 100644
--- a/main.go
+++ b/main.go
@@ -1,3 +1,5 @@
 package main

+import "fmt"
+
 func main() {
+	fmt.Println("hello world")
 }
`
	msg, err := generateCommitMessageWithClaude(diff)
	if err != nil {
		t.Fatalf("generateCommitMessageWithClaude failed: %v", err)
	}

	if msg == "" {
		t.Fatal("expected non-empty commit message")
	}

	if len(msg) > 100 {
		t.Errorf("commit message too long (%d chars): %s", len(msg), msg)
	}

	t.Logf("Generated commit message: %q", msg)
}

func TestGetCurrentBranch(t *testing.T) {
	dir, err := os.MkdirTemp("", "prepare-commit-msg.test")
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
