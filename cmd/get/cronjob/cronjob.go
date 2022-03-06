package cronjob

import (
	"errors"
	"fmt"

	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
)

// ErrResourceNameRequired is returned when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")

// ErrUnsupportedGroup is returned when a group/resource combination is invalid.
var ErrUnsupportedGroup = errors.New("group/resource not supported")

func clientError(err error) error {
	return fmt.Errorf("client error: %w", err)
}

func newRunError(err error) error {
	return fmt.Errorf("cronjob error: %w", err)
}

// NewCmd creates the `cronjob` command.
func NewCmd(opt *options.Options) *cobra.Command {
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

func validArgs(opt *options.Options) []string {
	client := client.NewClient(opt.Client, client.WithNamespace(opt.Namespace))
	group, _ := client.GetAPIGroup("CronJob")

	var list []string

	switch group {
	case "batch/v1beta1":
		list, _ = client.CronJobsV1beta1()
	case "batch/v1":
		list, _ = client.CronJobsV1()
	}

	return list
}

func run(opt *options.Options, args []string) error {
	if len(args) == 0 {
		return ErrResourceNameRequired
	}

	cl := client.NewClient(
		opt.Client,
		client.WithNamespace(opt.Namespace),
		client.WithFilename(opt.Filename),
		client.WithWriter(opt.Writer),
		client.WithExport(!opt.NoExport),
	)

	group, err := cl.GetAPIGroup("CronJob")
	if err != nil {
		return clientError(err)
	}

	switch group {
	case "batch/v1beta1":
		err = cl.CronJobV1Beta1(args[0]).Write()
	case "batch/v1":
		err = cl.CronJobV1(args[0]).Write()
	default:
		return ErrUnsupportedGroup
	}

	if err != nil {
		return newRunError(err)
	}

	return nil
}
