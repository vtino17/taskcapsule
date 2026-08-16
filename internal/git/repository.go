package git

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"strings"
)

func Root() (string, error) {
	root, err := execGit("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return canonicalRoot(root)
}

func RepoID(root string) (string, error) {
	canonical, err := canonicalRoot(root)
	if err != nil {
		return "", fmt.Errorf("canonicalize repository root: %w", err)
	}

	// Normalize path separators so that equivalent platform spellings produce
	// the same fallback ID. Resolving symlinks also handles macOS's /var ->
	// /private/var alias consistently.
	normalized := filepath.ToSlash(canonical)

	remote, err := execGitInDir(canonical, "remote", "get-url", "origin")
	if err != nil {
		h := sha256.Sum256([]byte(normalized))
		return fmt.Sprintf("%x", h[:8]), nil
	}

	h := sha256.Sum256([]byte(remote))
	return fmt.Sprintf("%x", h[:8]), nil
}

func canonicalRoot(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(filepath.Clean(resolved)), nil
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
