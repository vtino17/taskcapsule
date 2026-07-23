package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func init() {
	register(command{name: "completion", desc: "Generate shell completion scripts", handler: handleCompletion})
}

func commandNames() []string {
	var names []string
	for _, c := range commands {
		if c.name == "completion" {
			continue
		}
		names = append(names, c.name)
	}
	sort.Strings(names)
	return names
}

func handleCompletion(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: taskcapsule completion <shell>")
		fmt.Fprintln(os.Stderr, "Supported shells: bash, zsh, fish, powershell")
		return 2
	}

	shell := args[0]
	cmds := commandNames()

	switch shell {
	case "bash":
		completionBash(cmds)
		return 0
	case "zsh":
		completionZsh(cmds)
		return 0
	case "fish":
		completionFish(cmds)
		return 0
	case "powershell":
		completionPowerShell(cmds)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown shell %q\nSupported: bash, zsh, fish, powershell\n", shell)
		return 2
	}
}

func completionBash(cmds []string) {
	opts := strings.Join(append(cmds, "--help", "--version"), " ")
	fmt.Printf(`_taskcapsule() {
  local cur prev opts
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  opts="%s"
  COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
  return 0
}
complete -F _taskcapsule taskcapsule
`, opts)
}

func completionZsh(cmds []string) {
	opts := strings.Join(cmds, " ")
	fmt.Printf(`#compdef taskcapsule
_taskcapsule() {
  local -a opts
  opts=(%s)
  _describe 'taskcapsule' opts
}
_taskcapsule "$@"
`, opts)
}

func completionFish(cmds []string) {
	for _, cmd := range cmds {
		fmt.Printf("complete -c taskcapsule -f -a \"%s\"\n", cmd)
	}
}

func completionPowerShell(cmds []string) {
	opts := strings.Join(cmds, `", "`)
	fmt.Printf(`Register-ArgumentCompleter -Native -CommandName taskcapsule -ScriptBlock {
  param($wordToComplete, $commandAst, $cursorPosition)
  $commands = @("%s")
  $commands | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object { $_ }
}
`, opts)
}

func completionHelp() string {
	return "Usage: taskcapsule completion <bash|zsh|fish|powershell>\n\nGenerate shell completion scripts.\n\nExamples:\n  taskcapsule completion bash > ~/.local/share/bash-completion/completions/taskcapsule\n  taskcapsule completion zsh > ~/.zsh/completion/_taskcapsule\n  taskcapsule completion fish > ~/.config/fish/completions/taskcapsule.fish\n  taskcapsule completion powershell >> $PROFILE"
}
