package completion

import (
	"errors"
	"fmt"
	"io"

	"github.com/eiladin/k8s-dotenv/pkg/options"
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
To load completions:

Bash:

  $ source <(k8s-dotenv completion bash)

  # To load completions for each session, execute once:
  # Linux:
    $ k8s-dotenv completion bash > /etc/bash_completion.d/k8s-dotenv
  # macOS:
    $ k8s-dotenv completion bash > /usr/local/etc/bash_completion.d/k8s-dotenv
	
Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
	
    $ echo "autoload -U compinit; compinit" >> ~/.zshrc
	
  # To load completions for each session, execute once:
    $ k8s-dotenv completion zsh > "${fpath[1]}/_k8s-dotenv"
	
  # You will need to start a new shell for this setup to take effect.
	
fish:
	
  $ k8s-dotenv completion fish | source
	
  # To load completions for each session, execute once:
    $ k8s-dotenv completion fish > ~/.config/fish/completions/k8s-dotenv.fish
	
PowerShell:
	
  PS> k8s-dotenv completion powershell | Out-String | Invoke-Expression
	
  # To load completions for every new session, run:
    PS> k8s-dotenv completion powershell > k8s-dotenv.ps1
  # and source this file from your PowerShell profile.`
)

var (
	completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash":       runCompletionBash,
		"zsh":        runCompletionZsh,
		"fish":       runCompletionFish,
		"powershell": runCompletionPwsh,
	}
)

func NewCmd(opt *options.Options) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCompletion(opt, cmd, args)
		},
		ValidArgs: shells,
	}

	return cmd
}

func RunCompletion(opt *options.Options, cmd *cobra.Command, args []string) error {
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

	return run(opt.Writer, cmd.Parent())
}

func runCompletionBash(out io.Writer, root *cobra.Command) error {
	if _, err := out.Write([]byte(defaultBoilerPlate)); err != nil {
		return err
	}

	return root.GenBashCompletionV2(out, true)
}

func runCompletionZsh(out io.Writer, root *cobra.Command) error {
	zshHead := fmt.Sprintf("#compdef %[1]s\ncompdef _%[1]s %[1]s\n", root.Name())
	_, _ = out.Write([]byte(zshHead))

	if _, err := out.Write([]byte(defaultBoilerPlate)); err != nil {
		return err
	}

	return root.GenZshCompletion(out)
}

func runCompletionFish(out io.Writer, root *cobra.Command) error {
	if _, err := out.Write([]byte(defaultBoilerPlate)); err != nil {
		return err
	}

	return root.GenFishCompletion(out, true)
}

func runCompletionPwsh(out io.Writer, root *cobra.Command) error {

	if _, err := out.Write([]byte(defaultBoilerPlate)); err != nil {
		return err
	}

	return root.GenPowerShellCompletionWithDesc(out)
}
