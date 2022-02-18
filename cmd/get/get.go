package get

import (
	"github.com/eiladin/k8s-dotenv/cmd/get/cronjob"
	"github.com/eiladin/k8s-dotenv/cmd/get/daemonset"
	"github.com/eiladin/k8s-dotenv/cmd/get/deployment"
	"github.com/eiladin/k8s-dotenv/cmd/get/job"
	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/spf13/cobra"
)

func NewCmd(opt *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get DEPLOYMENT_NAME",
		Short: "fetch secrets and configmaps into a file",
	}

	cmd.AddCommand(cronjob.NewCmd(opt))
	cmd.AddCommand(deployment.NewCmd(opt))
	cmd.AddCommand(daemonset.NewCmd(opt))
	cmd.AddCommand(job.NewCmd(opt))

	return cmd
}
