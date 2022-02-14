package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const defaultBoilerPlate = `
# The MIT License (MIT)
# 
# Copyright © 2022 Sami Khan
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.
`

var (
	completionLong = `
Output shell completion code for the specified shell (bash, zsh, fish, or powershell). The shell code must be evaluated to provide interactive
completion of k8s-dotenv commands.  This can be done by sourcing it from the .bash_profile.
  Note for zsh users: zsh completions are only supported in versions of zsh >= 5.2.`

	completionExample = `
# Installing bash completion on macOS using homebrew
## If running Bash 3.2 included with macOS
  brew install bash-completion
## or, if running Bash 4.1+
  brew install bash-completion@2
## Add the completion to your completion directory
  k8s-dotenv completion bash > $(brew --prefix)/etc/bash_completion.d/k8s-dotenv

# Installing bash completion on Linux
## If bash-completion is not installed on Linux, install the 'bash-completion' package
## via your distribution's package manager.
## Load the k8s-dotenv completion code for bash into the current shell
  source <(k8s-dotenv completion bash)
## Write bash completion code to a file and source it from .bash_profile
  k8s-dotenv completion bash > ~/.k8s-dotenv/completion.bash.inc
  printf "
    # k8s-dotenv shell completion
    source '$HOME/.k8s-dotenv/completion.bash.inc'
    " >> $HOME/.bash_profile
  source $HOME/.bash_profile

	# Load the k8s-dotenv completion code for zsh into the current shell
  source <(k8s-dotenv completion zsh)
# Set the k8s-dotenv completion code for zsh to autoload on startup
  k8s-dotenv completion zsh > "${fpath[1]}/_k8s-dotenv"
# Load the k8s-dotenv completion code for fish into the current shell
  k8s-dotenv completion fish | source
# To load completions for each session, execute once: 
  k8s-dotenv completion fish > ~/.config/fish/completions/k8s-dotenv.fish
# Load the k8s-dotenv completion code for powershell into the current shell
  k8s-dotenv completion powershell | Out-String | Invoke-Expression
# Set k8s-dotenv completion code for powershell to run on startup
## Save completion code to a script and execute in the profile
  k8s-dotenv completion powershell > $HOME\.k8s-dotenv\completion.ps1
  Add-Content $PROFILE "$HOME\.k8s-dotenv\completion.ps1"
## Execute completion code in the profile
  Add-Content $PROFILE "if (Get-Command k8s-dotenv -ErrorAction SilentlyContinue) {
    k8s-dotenv completion powershell | Out-String | Invoke-Expression
  }"
## Add completion code directly to the $PROFILE script
  k8s-dotenv completion powershell >> $PROFILE`
)

var (
	completionShells = map[string]func(out io.Writer, boilerPlate string, cmd *cobra.Command) error{
		"bash":       runCompletionBash,
		"zsh":        runCompletionZsh,
		"fish":       runCompletionFish,
		"powershell": runCompletionPwsh,
	}
)

func newCompletionCmd(boilerPlate string) *cobra.Command {
	shells := []string{}
	for s := range completionShells {
		shells = append(shells, s)
	}

	cmd := &cobra.Command{
		Use:                   "completion SHELL",
		DisableFlagsInUseLine: true,
		Short:                 "Output shell completion code for the specified shell (bash, zsh, fish)",
		Long:                  completionLong,
		Example:               completionExample,
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunCompletion(os.Stdout, boilerPlate, cmd, args); err != nil {
				log.Fatal(err)
			}
		},
		ValidArgs: shells,
	}

	return cmd
}

func RunCompletion(out io.Writer, boilerPlate string, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("shell not specified")
	}
	if len(args) > 1 {
		return errors.New("Too many arguments. Expected only the shell type.")
	}
	run, found := completionShells[args[0]]
	if !found {
		return fmt.Errorf("Unsupported shell type %q.", args[0])
	}

	return run(out, boilerPlate, cmd.Parent())
}

func runCompletionBash(out io.Writer, boilerPlate string, root *cobra.Command) error {
	if len(boilerPlate) == 0 {
		boilerPlate = defaultBoilerPlate
	}
	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	return root.GenBashCompletionV2(out, false) // TODO: Upgrade to Cobra 1.3.0 or later before including descriptions (See https://github.com/spf13/cobra/pull/1509)
}

func runCompletionZsh(out io.Writer, boilerPlate string, root *cobra.Command) error {
	zshHead := fmt.Sprintf("#compdef %[1]s\ncompdef _%[1]s %[1]s\n", root.Name())
	_, _ = out.Write([]byte(zshHead))

	if len(boilerPlate) == 0 {
		boilerPlate = defaultBoilerPlate
	}
	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	return root.GenZshCompletion(out)
}

func runCompletionFish(out io.Writer, boilerPlate string, root *cobra.Command) error {
	if len(boilerPlate) == 0 {
		boilerPlate = defaultBoilerPlate
	}
	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	return root.GenFishCompletion(out, true)
}

func runCompletionPwsh(out io.Writer, boilerPlate string, root *cobra.Command) error {
	if len(boilerPlate) == 0 {
		boilerPlate = defaultBoilerPlate
	}

	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	return root.GenPowerShellCompletionWithDesc(out)
}
