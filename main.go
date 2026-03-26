package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	var dat []byte
	var err error
	dat, err = os.ReadFile(os.Args[1])
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
	err = os.WriteFile(os.Args[1], []byte(prepareMessage(branchName, message)), 0644)
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

type claudeStreamMessage struct {
	Type    string `json:"type"`
	Message struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"message"`
}

func generateCommitMessageWithClaude(diff string) (string, error) {
	prompt := fmt.Sprintf("Based on this git diff, write a very concise commit message (one line, maximum 50 characters). Only output the commit message itself, nothing else.\n\nGit diff:\n%s", diff)

	cmd := exec.Command("claude", "-p", "--model", "haiku", "--output-format", "stream-json", "--verbose", prompt)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, stderr.String())
	}

	// Parse stream-json output: each line is a JSON object.
	// Look for assistant messages with text content blocks.
	var textParts []string
	for _, line := range strings.Split(stdout.String(), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var msg claudeStreamMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}
		if msg.Type != "assistant" {
			continue
		}
		for _, block := range msg.Message.Content {
			if block.Type == "text" {
				textParts = append(textParts, block.Text)
			}
		}
	}

	combined := strings.TrimSpace(strings.Join(textParts, ""))
	if combined == "" {
		return "", fmt.Errorf("no text content in claude response")
	}

	// Strip markdown code fences if present
	lines := strings.Split(combined, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "```") {
			continue
		}
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
