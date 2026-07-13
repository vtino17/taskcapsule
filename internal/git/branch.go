package git

import (
	"strings"
)

func BranchExists(branch, repoRoot string) (bool, error) {
	_, err := execGitInDir(repoRoot, "rev-parse", "--verify", "refs/heads/"+branch)
	if err != nil {
		// Try remote branch
		_, err2 := execGitInDir(repoRoot, "rev-parse", "--verify", "refs/remotes/origin/"+branch)
		return err2 == nil, nil
	}
	return true, nil
}

func BranchInUse(branch, repoRoot string) (bool, error) {
	out, err := execGitInDir(repoRoot, "worktree", "list")
	if err != nil {
		return false, err
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "["+branch+"]") {
			return true, nil
		}
	}
	return false, nil
}

func CreateBranch(repoRoot, branch, baseBranch string) error {
	// Check if branch already exists
	exists, _ := BranchExists(branch, repoRoot)
	if exists {
		return nil
	}

	_, err := execGitInDir(repoRoot, "branch", branch, baseBranch)
	return err
}
