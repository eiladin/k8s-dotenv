package cronjob

import (
	"errors"
	"fmt"

	v1 "github.com/eiladin/k8s-dotenv/pkg/api/v1"
	"github.com/eiladin/k8s-dotenv/pkg/api/v1beta1"
	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
)

// ErrResourceNameRequired is returned when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")

// ErrUnsupportedGroup is returned when a group/resource combination is invalid.
var ErrUnsupportedGroup = errors.New("group/resource not supported")

func newClientError(err error) error {
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
	group, _ := opt.Client.GetAPIGroup("CronJob")

	var list []string

	switch group {
	case "batch/v1beta1":
		list, _ = v1beta1.CronJobs(opt.Client, opt.Namespace)
	case "batch/v1":
		list, _ = v1.CronJobs(opt.Client, opt.Namespace)
	}

	return list
}

func run(opt *options.Options, args []string) error {
	if len(args) == 0 {
		return ErrResourceNameRequired
	}

	group, err := opt.Client.GetAPIGroup("CronJob")
	if err != nil {
		return newClientError(err)
	}

	var res *environment.Result

	switch group {
	case "batch/v1beta1":
		res, err = v1beta1.CronJob(opt.Client, opt.Namespace, args[0])
	case "batch/v1":
		res, err = v1.CronJob(opt.Client, opt.Namespace, args[0])
	default:
		return ErrUnsupportedGroup
	}

	if err != nil {
		return newRunError(err)
	}

	if err := res.Write(opt); err != nil {
		return newRunError(err)
	}

	return nil
}
