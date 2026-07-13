package report

import (
	"fmt"
	"strings"
)

type HandoffData struct {
	Name            string
	Status          string
	Branch          string
	BaseBranch      string
	Dirty           bool
	ChangedFiles    []string
	CurrentNote     string
	Services        []string
	LastCheckCmd    string
	LastCheckResult string
	LastCheckExit   int
}

func GenerateHandoff(data HandoffData) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# Handoff: %s\n\n", data.Name))

	b.WriteString("## Status\n\n")
	b.WriteString(data.Status + "\n\n")

	if data.CurrentNote != "" {
		b.WriteString("## Current objective\n\n")
		b.WriteString(data.CurrentNote + "\n\n")
	}

	b.WriteString("## Git\n\n")
	b.WriteString(fmt.Sprintf("- Branch: %s\n", data.Branch))
	b.WriteString(fmt.Sprintf("- Base: %s\n", data.BaseBranch))
	dirtyStr := "no"
	if data.Dirty {
		dirtyStr = "yes"
	}
	b.WriteString(fmt.Sprintf("- Dirty: %s\n\n", dirtyStr))

	if len(data.ChangedFiles) > 0 {
		b.WriteString("## Changed files\n\n")
		for _, f := range data.ChangedFiles {
			b.WriteString(fmt.Sprintf("- %s\n", f))
		}
		b.WriteString("\n")
	}

	if data.LastCheckCmd != "" {
		b.WriteString("## Last check\n\n")
		b.WriteString(fmt.Sprintf("Command: `%s`\n", data.LastCheckCmd))
		b.WriteString(fmt.Sprintf("Result: %s\n", data.LastCheckResult))
		b.WriteString(fmt.Sprintf("Exit code: %d\n\n", data.LastCheckExit))
	}

	if len(data.Services) > 0 {
		b.WriteString("## Services\n\n")
		for _, svc := range data.Services {
			b.WriteString(fmt.Sprintf("- %s\n", svc))
		}
		b.WriteString("\n")
	}

	b.WriteString("## How to continue\n\n")
	b.WriteString("```bash\n")
	b.WriteString(fmt.Sprintf("taskcapsule resume %s\n", data.Name))
	b.WriteString("```\n\n")

	b.WriteString("## Security\n\n")
	b.WriteString("No environment variable values are included.\n")

	return b.String()
}
