package app

import (
	"fmt"

	"github.com/vtino17/taskcapsule/internal/lock"
	"github.com/vtino17/taskcapsule/internal/state"
)

type capsuleLock struct {
	lock *lock.Lock
	repo string
	name string
}

func acquireCapsuleLock(repoID, capsuleName, command string) (*capsuleLock, error) {
	stateBase, err := getStateDir()
	if err != nil {
		return nil, err
	}

	mgr := lock.NewManager(stateBase)
	l, err := mgr.Acquire(repoID, capsuleName, command)
	if err != nil {
		return nil, fmt.Errorf("%v\nExit code: 4", err)
	}

	return &capsuleLock{lock: l, repo: repoID, name: capsuleName}, nil
}

func (cl *capsuleLock) Release() {
	if cl == nil || cl.lock == nil {
		return
	}
	cl.lock.Release()
}

func lockAndLoad(repoID, capsuleName, command string) (*capsuleLock, *state.Store, error) {
	cl, err := acquireCapsuleLock(repoID, capsuleName, command)
	if err != nil {
		return nil, nil, err
	}

	stateBase, _ := getStateDir()
	cs := state.NewStore(stateBase)
	s, err := cs.Load(repoID, capsuleName)
	if err != nil {
		cl.Release()
		return nil, nil, err
	}
	_ = s // Caller uses returned store

	return cl, cs, nil
}
