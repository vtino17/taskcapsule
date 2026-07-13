package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateWorktree(repoRoot, worktreePath, branch, baseBranch string) error {
	// Check if worktree already exists
	exists, _ := worktreeExists(repoRoot, worktreePath)
	if exists {
		return nil
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(worktreePath), 0755); err != nil {
		return fmt.Errorf("cannot create worktree parent dir: %v", err)
	}

	// Create branch if needed
	if err := CreateBranch(repoRoot, branch, baseBranch); err != nil {
		return err
	}

	_, err := execGitInDir(repoRoot, "worktree", "add", worktreePath, branch)
	if err != nil {
		// Try with -b flag if branch doesn't exist yet
		exists, _ := BranchExists(branch, repoRoot)
		if !exists {
			_, err = execGitInDir(repoRoot, "worktree", "add", "-b", branch, worktreePath, baseBranch)
		}
	}
	return err
}

func RemoveWorktree(worktreePath string) error {
	_, err := execGit("worktree", "remove", worktreePath)
	if err != nil {
		// Try with --force
		_, err2 := execGit("worktree", "remove", "--force", worktreePath)
		return err2
	}
	return err
}

func IsDirty(worktreePath string) (bool, error) {
	out, err := execGitInDir(worktreePath, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

func ChangedFiles(worktreePath string) ([]string, error) {
	out, err := execGitInDir(worktreePath, "status", "--porcelain")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	var files []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Format: "XY filename" or "XY  filename"
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			files = append(files, parts[1])
		}
	}
	return files, nil
}

func worktreeExists(repoRoot, worktreePath string) (bool, error) {
	out, err := execGitInDir(repoRoot, "worktree", "list")
	if err != nil {
		return false, err
	}
	return strings.Contains(out, worktreePath), nil
}

func FindWorktreeUsingBranch(repoRoot, branch string) (string, error) {
	out, err := execGitInDir(repoRoot, "worktree", "list")
	if err != nil {
		return "", err
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "["+branch+"]") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				return parts[0], nil
			}
		}
	}
	return "", fmt.Errorf("no worktree found for branch %q", branch)
}
