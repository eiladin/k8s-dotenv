package get

import (
	"github.com/eiladin/k8s-dotenv/cmd/get/cronjob"
	"github.com/eiladin/k8s-dotenv/cmd/get/daemonset"
	"github.com/eiladin/k8s-dotenv/cmd/get/deployment"
	"github.com/eiladin/k8s-dotenv/cmd/get/job"
	"github.com/eiladin/k8s-dotenv/cmd/get/pod"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
)

// NewCmd creates the `get` command.
func NewCmd(opt *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get RESOURCE_TYPE",
		Short: "fetch secrets and configmaps into a file",
	}

	cmd.AddCommand(cronjob.NewCmd(opt))
	cmd.AddCommand(deployment.NewCmd(opt))
	cmd.AddCommand(daemonset.NewCmd(opt))
	cmd.AddCommand(job.NewCmd(opt))
	cmd.AddCommand(pod.NewCmd(opt))

	return cmd
}
