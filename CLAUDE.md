# prepare-commit-msg

Git hook that prepends Jira ticket IDs from branch names to commit messages. When the commit message is "wip", it uses the Claude CLI to generate a meaningful message from the staged diff.

## Testing

Run unit tests:
```
go test -short ./...
```

Run all tests including integration tests:
```
go test -v ./...
```

**IMPORTANT:** The integration test (`TestGenerateCommitMessageWithClaude`) must pass before creating a new release. This test calls the Claude CLI end-to-end and verifies that commit message generation actually works. Do not tag a release without running `go test -v ./...` first.

## Releases

Releases are triggered by pushing a git tag. To create a new release:

1. Ensure all tests pass: `go test -v ./...`
2. Tag the commit: `git tag vX.X.X`
3. Push the tag: `git push origin vX.X.X`

There is no release script or CI-based tagging -- tags must be created and pushed manually.
