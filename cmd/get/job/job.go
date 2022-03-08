package job

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
	return fmt.Errorf("job error: %w", err)
}

// NewCmd creates the `job` command.
func NewCmd(opt *clioptions.CLIOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "job RESOURCE_NAME",
		Aliases: []string{"jobs"},
		Short:   "fetch environment configuration from job into a file",
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
	).BatchV1().JobList()

	return list
}

func run(opt *clioptions.CLIOptions, args []string) error {
	if len(args) == 0 {
		return ErrResourceNameRequired
	}

	err := client.NewClient(
		client.WithKubeClient(opt.KubeClient),
		client.WithNamespace(opt.Namespace),
		client.WithExport(!opt.NoExport),
	).BatchV1().Job(args[0]).Write(opt.Writer)

	if err != nil {
		return runError(err)
	}

	return nil
}
