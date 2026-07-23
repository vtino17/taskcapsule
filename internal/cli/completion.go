package cli

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

func init() {
	register(command{name: "completion", desc: "Generate shell completion scripts", handler: handleCompletion})
}

func handleCompletion(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: taskcapsule completion <shell>")
		fmt.Fprintln(os.Stderr, "Supported shells: bash, zsh, fish, powershell")
		return 2
	}

	shell := args[0]
	switch shell {
	case "bash":
		completionBash()
		return 0
	case "zsh":
		completionZsh()
		return 0
	case "fish":
		completionFish()
		return 0
	case "powershell":
		completionPowerShell()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "Unknown shell: %s\nSupported: bash, zsh, fish, powershell\n", shell)
		return 2
	}
}

func commandNames() []string {
	var names []string
	for _, c := range commands {
		names = append(names, c.name)
	}
	sort.Strings(names)
	return names
}

func completionBash() {
	cmds := strings.Join(commandNames(), " ")
	fmt.Println(`_taskcapsule() {
  local cur prev opts
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  opts="` + cmds + ` --help --version"
  COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
  return 0
}
complete -F _taskcapsule taskcapsule`)
}

func completionZsh() {
	cmds := strings.Join(commandNames(), " ")
	fmt.Println(`#compdef taskcapsule
_taskcapsule() {
  local -a opts
  opts=(` + cmds + `)
  _describe 'taskcapsule' opts
}
_taskcapsule "$@"`)
}

func completionFish() {
	cmds := strings.Join(commandNames(), " ")
	fmt.Println(`complete -c taskcapsule -f -a "` + cmds + `"`)
}

func completionPowerShell() {
	fmt.Println(`@("completion", "check", "delete", "doctor", "handoff", "init", "list", "logs", "note", "pause", "resume", "start", "status", "version", "where") | ForEach-Object {
  Register-ArgumentCompleter -Native -CommandName taskcapsule -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $commands = @("completion", "check", "delete", "doctor", "handoff", "init", "list", "logs", "note", "pause", "resume", "start", "status", "version", "where")
    $commands | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object { $_ }
  }
}`)
}

func CompletionUsage() string {
	return "Usage: taskcapsule completion <bash|zsh|fish|powershell>\n\nGenerate shell completion scripts.\n\nExample:\n  taskcapsule completion bash > /etc/bash_completion.d/taskcapsule"
}

var ErrCompletionUsage = errors.New(CompletionUsage())
