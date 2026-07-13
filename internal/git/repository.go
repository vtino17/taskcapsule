package git

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func Root() (string, error) {
	return execGit("rev-parse", "--show-toplevel")
}

func RepoID(root string) (string, error) {
	remote, err := execGitInDir(root, "remote", "get-url", "origin")
	if err != nil {
		// Fallback: use repo root path hash
		h := sha256.Sum256([]byte(root))
		return fmt.Sprintf("%x", h[:8]), nil
	}

	h := sha256.Sum256([]byte(remote))
	return fmt.Sprintf("%x", h[:8]), nil
}

func RepoName(root string) (string, error) {
	remote, err := execGitInDir(root, "remote", "get-url", "origin")
	if err != nil {
		// Fallback: use directory name
		parts := strings.Split(strings.TrimRight(root, "/"), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
		return "repo", nil
	}

	// Extract repo name from remote URL
	remote = strings.TrimSpace(remote)
	remote = strings.TrimSuffix(remote, ".git")

	parts := strings.Split(remote, "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		name = strings.TrimSuffix(name, ".git")
		return name, nil
	}

	return "repo", nil
}

func CurrentBranch(root string) (string, error) {
	return execGitInDir(root, "rev-parse", "--abbrev-ref", "HEAD")
}

func DefaultBranch(root string) (string, error) {
	branch, err := execGitInDir(root, "symbolic-ref", "refs/remotes/origin/HEAD")
	if err != nil {
		return "main", nil
	}
	branch = strings.TrimPrefix(branch, "refs/remotes/origin/")
	return strings.TrimSpace(branch), nil
}
