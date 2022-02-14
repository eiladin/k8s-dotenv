package cmd

import (
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get DEPLOYMENT_NAME",
		Short: "fetch secrets and configmaps into a file",
	}

	cmd.AddCommand(newCronJobCmd())
	cmd.AddCommand(newDeployCmd())
	cmd.AddCommand(newDaemonSetCmd())
	cmd.AddCommand(newJobCmd())

	return cmd
}
