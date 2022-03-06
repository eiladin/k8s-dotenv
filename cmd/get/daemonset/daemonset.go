package daemonset

import (
	"errors"
	"fmt"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/clioptions"
	"github.com/spf13/cobra"
)

// ErrResourceNameRequired is returned when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")

func runError(err error) error {
	return fmt.Errorf("daemonset error: %w", err)
}

// NewCmd creates the `daemonset` command.
func NewCmd(opt *clioptions.CLIOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "daemonset RESOURCE_NAME",
		Aliases: []string{"daemonsets", "ds"},
		Short:   "fetch environment configuration from daemon set into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return validArgs(opt), cobra.ShellCompDirectiveDefault
		},
		RunE: func(c *cobra.Command, args []string) error {
			return run(opt, args)
		},
	}

	return cmd
}

func validArgs(opt *clioptions.CLIOptions) []string {
	list, _ := client.NewClient(
		client.WithKubeClient(opt.KubeClient),
		client.WithNamespace(opt.Namespace),
	).AppsV1().DaemonSets()

	return list
}

func run(opt *clioptions.CLIOptions, args []string) error {
	if len(args) == 0 {
		return ErrResourceNameRequired
	}

	err := client.NewClient(
		client.WithKubeClient(opt.KubeClient),
		client.WithNamespace(opt.Namespace),
		client.WithFilename(opt.Filename),
		client.WithWriter(opt.Writer),
		client.WithExport(!opt.NoExport),
	).AppsV1().DaemonSet(args[0]).Write()

	if err != nil {
		return runError(err)
	}

	return nil
}
