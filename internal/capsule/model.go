package capsule

import "time"

type State struct {
	SchemaVersion  int                     `json:"schemaVersion"`
	Name           string                  `json:"name"`
	Status         string                  `json:"status"`
	RepositoryRoot string                  `json:"repositoryRoot"`
	RepositoryID   string                  `json:"repositoryID"`
	WorktreePath   string                  `json:"worktreePath"`
	Branch         string                  `json:"branch"`
	BaseBranch     string                  `json:"baseBranch"`
	CreatedAt      time.Time               `json:"createdAt"`
	UpdatedAt      time.Time               `json:"updatedAt"`
	LastPausedAt   *time.Time              `json:"lastPausedAt,omitempty"`
	CurrentNote    string                  `json:"currentNote"`
	NoteHistory    []NoteEntry             `json:"noteHistory,omitempty"`
	Services       map[string]ServiceState `json:"services,omitempty"`
	LastCheck      *CheckState             `json:"lastCheck,omitempty"`
	LastError      *string                 `json:"lastError,omitempty"`
}

type NoteEntry struct {
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

type ServiceState struct {
	Status        string   `json:"status"`
	Command       []string `json:"command"`
	PID           int      `json:"pid"`
	ProcessGroup  int      `json:"processGroupID"`
	Port          int      `json:"port"`
	LogPath       string   `json:"logPath"`
	LastStartedAt string   `json:"lastStartedAt"`
	LastStoppedAt string   `json:"lastStoppedAt"`
}

type CheckState struct {
	Command    []string `json:"command"`
	ExitCode   int      `json:"exitCode"`
	StartedAt  string   `json:"startedAt"`
	FinishedAt string   `json:"finishedAt"`
	LogPath    string   `json:"logPath"`
}
