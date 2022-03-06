package daemonset

import (
	"errors"
	"fmt"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
)

// ErrResourceNameRequired is returned when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")

func newRunError(err error) error {
	return fmt.Errorf("daemonset error: %w", err)
}

// NewCmd creates the `daemonset` command.
func NewCmd(opt *options.Options) *cobra.Command {
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

func validArgs(opt *options.Options) []string {
	list, _ := client.NewClient(
		opt.Client,
		client.WithNamespace(opt.Namespace),
	).DaemonSetsV1()

	return list
}

func run(opt *options.Options, args []string) error {
	if len(args) == 0 {
		return ErrResourceNameRequired
	}

	err := client.NewClient(
		opt.Client,
		client.WithNamespace(opt.Namespace),
		client.WithFilename(opt.Filename),
		client.WithWriter(opt.Writer),
		client.WithExport(!opt.NoExport),
	).DaemonSetV1(args[0]).Write()

	if err != nil {
		return newRunError(err)
	}

	return nil
}
