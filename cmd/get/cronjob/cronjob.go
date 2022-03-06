package cronjob

import (
	"errors"
	"fmt"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/clioptions"
	"github.com/spf13/cobra"
)

// ErrResourceNameRequired is returned when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")

// ErrUnsupportedGroup is returned when a group/resource combination is invalid.
var ErrUnsupportedGroup = errors.New("group/resource not supported")

func clientError(err error) error {
	return fmt.Errorf("client error: %w", err)
}

func runError(err error) error {
	return fmt.Errorf("cronjob error: %w", err)
}

// NewCmd creates the `cronjob` command.
func NewCmd(opt *clioptions.CLIOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cronjob RESOURCE_NAME",
		Aliases: []string{"cronjobs", "cj"},
		Short:   "fetch environment configuration from cron job into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return validArgs(opt), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(c *cobra.Command, args []string) error {
			return run(opt, args)
		},
	}

	return cmd
}

func validArgs(opt *clioptions.CLIOptions) []string {
	client := client.NewClient(
		client.WithKubeClient(opt.KubeClient),
		client.WithNamespace(opt.Namespace),
	)
	group, _ := client.GetAPIGroup("CronJob")

	var list []string

	switch group {
	case "batch/v1beta1":
		list, _ = client.BatchV1Beta1().CronJobs()
	case "batch/v1":
		list, _ = client.BatchV1().CronJobs()
	}

	return list
}

func run(opt *clioptions.CLIOptions, args []string) error {
	if len(args) == 0 {
		return ErrResourceNameRequired
	}

	client := client.NewClient(
		client.WithKubeClient(opt.KubeClient),
		client.WithNamespace(opt.Namespace),
		client.WithFilename(opt.Filename),
		client.WithWriter(opt.Writer),
		client.WithExport(!opt.NoExport),
	)

	group, err := client.GetAPIGroup("CronJob")
	if err != nil {
		return clientError(err)
	}

	switch group {
	case "batch/v1beta1":
		err = client.BatchV1Beta1().CronJob(args[0]).Write()
	case "batch/v1":
		err = client.BatchV1().CronJob(args[0]).Write()
	default:
		return ErrUnsupportedGroup
	}

	if err != nil {
		return runError(err)
	}

	return nil
}
