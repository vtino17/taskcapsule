package app

import (
	"fmt"
	"time"

	"github.com/vtino17/taskcapsule/internal/capsule"
	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/state"
)

func SaveNote(name, text string) error {
	if err := validateCapsuleName(name); err != nil {
		return err
	}

	root, err := findGitRoot()
	if err != nil {
		return err
	}

	repoID, err := git.RepoID(root)
	if err != nil {
		return err
	}

	cl, err := acquireCapsuleLock(repoID, name, "note")
	if err != nil {
		return err
	}
	defer cl.Release()

	stateBase, err := getStateDir()
	if err != nil {
		return err
	}
	cs := state.NewStore(stateBase)

	s, err := loadValidatedCapsule(cs, stateBase, repoID, name)
	if err != nil {
		return fmt.Errorf("capsule not found: %s", name)
	}

	if s.CurrentNote != "" {
		s.NoteHistory = append(s.NoteHistory, capsule.NoteEntry{
			Text:      s.CurrentNote,
			Timestamp: s.UpdatedAt.Format(time.RFC3339),
		})
	}

	s.CurrentNote = text
	s.UpdatedAt = time.Now().UTC()

	return cs.Save(repoID, name, s)
}
