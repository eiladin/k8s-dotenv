package cronjob

import (
	"os"

	v1 "github.com/eiladin/k8s-dotenv/internal/api/v1"
	"github.com/eiladin/k8s-dotenv/internal/api/v1beta1"
	"github.com/eiladin/k8s-dotenv/internal/client"
	"github.com/eiladin/k8s-dotenv/internal/environment"
	"github.com/eiladin/k8s-dotenv/internal/errors/cmd"
	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/spf13/cobra"
)

func NewCmd(opt *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cronjob RESOURCE_NAME",
		Aliases: []string{"cronjobs", "cj"},
		Short:   "fetch environment configuration from cron job into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			group, _ := client.GetApiGroup(opt.Client, "CronJob")

			var list []string
			if group == "batch/v1beta1" {
				list, _ = v1beta1.CronJobs(opt)
			} else {
				list, _ = v1.CronJobs(opt)
			}

			return list, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.ErrResourceNameRequired
			}
			group, err := client.GetApiGroup(opt.Client, "CronJob")
			if err != nil {
				return err
			}
			beta1 := group == "batch/v1beta1"

			var res *environment.Result
			opt.Name = args[0]
			if beta1 {
				res, err = v1beta1.CronJob(opt)
			} else {
				res, err = v1.CronJob(opt)
			}
			if err != nil {
				return err
			}

			f, err := os.OpenFile(opt.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				return err
			}
			defer f.Close()
			return res.Write(f, opt)
		},
	}

	return cmd
}
