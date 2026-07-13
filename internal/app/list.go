package app

import (
	"fmt"
	"sort"
	"time"

	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/state"
)

func List(showAll bool) ([]CapsuleInfo, error) {
	stateBase, err := getStateDir()
	if err != nil {
		return nil, err
	}

	// If showAll, list all capsules from all repos
	if showAll {
		return listAllCapsules(stateBase)
	}

	// Otherwise only list capsules in current repo
	root, err := findGitRoot()
	if err != nil {
		return nil, err
	}

	repoID, err := git.RepoID(root)
	if err != nil {
		return nil, err
	}

	cs := state.NewStore(stateBase)
	capsules, err := cs.List(repoID)
	if err != nil {
		return nil, err
	}

	result := make([]CapsuleInfo, 0, len(capsules))
	for _, c := range capsules {
		updated := time.Since(c.UpdatedAt).Truncate(time.Minute)
		updatedStr := updated.String()
		if updated < time.Minute {
			updatedStr = "just now"
		} else if updated < time.Hour {
			updatedStr = formatDuration(updated)
		} else {
			updatedStr = formatDuration(updated)
		}

		result = append(result, CapsuleInfo{
			Name:    c.Name,
			Status:  c.Status,
			Branch:  c.Branch,
			Updated: updatedStr,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func listAllCapsules(stateBase string) ([]CapsuleInfo, error) {
	cs := state.NewStore(stateBase)
	capsules, err := cs.ListAll()
	if err != nil {
		return nil, err
	}

	result := make([]CapsuleInfo, 0, len(capsules))
	for _, c := range capsules {
		updated := time.Since(c.UpdatedAt).Truncate(time.Minute)
		result = append(result, CapsuleInfo{
			Name:    c.Name,
			Status:  c.Status,
			Branch:  c.Branch,
			Updated: formatDuration(updated),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60

	if h > 24 {
		days := h / 24
		return fmt.Sprintf("%dd ago", days)
	}
	if h > 0 {
		return fmt.Sprintf("%dh ago", h)
	}
	if m > 0 {
		return fmt.Sprintf("%dm ago", m)
	}
	return "just now"
}
