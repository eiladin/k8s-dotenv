package cmd

import (
	v1 "github.com/eiladin/k8s-dotenv/internal/api/v1"
	"github.com/eiladin/k8s-dotenv/internal/api/v1beta1"
	"github.com/eiladin/k8s-dotenv/internal/environment"
	"github.com/spf13/cobra"
)

func newCronJobCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cronjob RESOURCE_NAME",
		Aliases: []string{"cronjobs", "cj"},
		Short:   "fetch environment configuration from cron job into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			group, _ := opt.Client.GetApiGroup("CronJob")

			var list []string
			if group == "batch/v1beta1" {
				list, _ = v1beta1.CronJobs(opt)
			} else {
				list, _ = v1.CronJobs(opt)
			}

			return list, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return ErrResourceNameRequired
			}
			group, err := opt.Client.GetApiGroup("CronJob")
			if err != nil {
				return err
			}
			beta1 := group == "batch/v1beta1"

			var envRes *environment.Result
			opt.Name = args[0]
			if beta1 {
				envRes, err = v1beta1.CronJob(opt)
			} else {
				envRes, err = v1.CronJob(opt)
			}
			if err != nil {
				return err
			}

			return envRes.Write(opt)
		},
	}

	return cmd
}
