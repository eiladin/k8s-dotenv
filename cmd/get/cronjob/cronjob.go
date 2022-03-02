package cronjob

import (
	"fmt"

	v1 "github.com/eiladin/k8s-dotenv/pkg/api/v1"
	"github.com/eiladin/k8s-dotenv/pkg/api/v1beta1"
	"github.com/eiladin/k8s-dotenv/pkg/environment"
	"github.com/eiladin/k8s-dotenv/pkg/errors/cmd"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
)

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
	if group == "batch/v1beta1" {
		list, _ = v1beta1.CronJobs(opt.Client, opt.Namespace)
	} else if group == "batch/v1" {
		list, _ = v1.CronJobs(opt.Client, opt.Namespace)
	} else {
		return list
	}

	return list
}

func run(opt *options.Options, args []string) error {
	if len(args) == 0 {
		return cmd.ErrResourceNameRequired
	}
	group, err := opt.Client.GetAPIGroup("CronJob")
	if err != nil {
		return err
	}

	var res *environment.Result
	if group == "batch/v1beta1" {
		res, err = v1beta1.CronJob(opt.Client, opt.Namespace, args[0])
	} else if group == "batch/v1" {
		res, err = v1.CronJob(opt.Client, opt.Namespace, args[0])
	} else {
		return fmt.Errorf("resource CronJob in group %s not supported", group)
	}
	if err != nil {
		return err
	}
	return res.Write(opt)
}
