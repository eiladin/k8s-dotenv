package pod

import (
	"errors"
	"fmt"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
)

// ErrResourceNameRequired is returned when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")

func runError(err error) error {
	return fmt.Errorf("pod error: %w", err)
}

// NewCmd creates the `pod` command.
func NewCmd(opt *options.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pod RESOURCE_NAME",
		Aliases: []string{"pods", "po"},
		Short:   "fetch environment configuration from pod into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return validArgs(opt), cobra.ShellCompDirectiveDefault
		},
		RunE: func(c *cobra.Command, args []string) error {
			return run(opt, args)
		},
	}

	return cmd
}

func validArgs(opt *options.CLI) []string {
	list, _ := client.NewClient(
		client.WithKubeClient(opt.KubeClient),
		client.WithNamespace(opt.Namespace),
	).CoreV1().PodList()

	return list
}

func run(opt *options.CLI, args []string) error {
	if len(args) == 0 {
		return ErrResourceNameRequired
	}

	err := client.NewClient(
		client.WithKubeClient(opt.KubeClient),
		client.WithNamespace(opt.Namespace),
		client.WithExport(!opt.NoExport),
	).CoreV1().Pod(args[0]).Write(opt.Writer)

	if err != nil {
		return runError(err)
	}

	return nil
}
