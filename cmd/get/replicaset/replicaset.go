package replicaset

import (
	"errors"
	"fmt"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/clioptions"
	"github.com/spf13/cobra"
)

// ErrResourceNameRequired is returned when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")

func newRunError(err error) error {
	return fmt.Errorf("replicaset error: %w", err)
}

// NewCmd creates the `replicaset` command.
func NewCmd(opt *clioptions.CLIOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "replicaset RESOURCE_NAME",
		Aliases: []string{"replicasets", "rs"},
		Short:   "fetch environment configuration from replica set into a file",
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
		opt.KubeClient,
		client.WithNamespace(opt.Namespace),
	).AppsV1().ReplicaSets()

	return list
}

func run(opt *clioptions.CLIOptions, args []string) error {
	if len(args) == 0 {
		return ErrResourceNameRequired
	}

	err := client.NewClient(
		opt.KubeClient,
		client.WithNamespace(opt.Namespace),
		client.WithFilename(opt.Filename),
		client.WithWriter(opt.Writer),
		client.WithExport(!opt.NoExport),
	).AppsV1().ReplicaSet(args[0]).Write()

	if err != nil {
		return newRunError(err)
	}

	return nil
}
