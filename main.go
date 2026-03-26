package main

import (
	"bytes"
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

	message := string(dat)

	// Check if the message is 'wip' and replace it with Claude-generated message
	firstLine := strings.Split(message, "\n")[0]
	if isWipMessage(firstLine) {
		fmt.Fprintln(os.Stderr, "🤖 Detected 'wip' commit message. Generating a better message with Claude...")

		diff, err := getStagedDiff()
		if err != nil || diff == "" {
			fmt.Fprintln(os.Stderr, "⚠️  No staged changes found. Keeping default message.")
		} else {
			newMsg, err := generateCommitMessageWithClaude(diff)
			if err != nil || newMsg == "" {
				fmt.Fprintln(os.Stderr, "❌ Failed to generate message with Claude. Keeping 'wip'.")
				fmt.Fprintln(os.Stderr, "   Make sure 'claude' CLI is installed and accessible.")
			} else {
				message = newMsg
				fmt.Fprintf(os.Stderr, "✅ Updated commit message to: %s\n", newMsg)
			}
		}
	}

	branchName, err := getCurrentBranch()
	if err != nil {
		return
	}
	err = ioutil.WriteFile(os.Args[1], []byte(prepareMessage(branchName, message)), 0644)
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

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func generateCommitMessageWithClaude(diff string) (string, error) {
	prompt := fmt.Sprintf("Based on this git diff, write a very concise commit message (one line, maximum 50 characters). Only output the commit message itself, nothing else.\n\nGit diff:\n%s", diff)

	cmd := exec.Command("claude", "-p", "--model", "haiku", prompt)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Get the output and strip markdown code fences if present
	output := strings.TrimSpace(stdout.String())
	lines := strings.Split(output, "\n")

	// Strip markdown code fences and find the actual message
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Skip empty lines and code fence markers
		if trimmedLine == "" || trimmedLine == "```" || strings.HasPrefix(trimmedLine, "```") {
			continue
		}
		// Return the first non-fence line
		return trimmedLine, nil
	}

	return "", fmt.Errorf("no valid commit message found in output")
}

func isWipMessage(message string) bool {
	trimmed := strings.TrimSpace(message)
	return strings.ToLower(trimmed) == "wip"
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
