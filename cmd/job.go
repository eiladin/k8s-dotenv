package cmd

import (
	"log"

	v1 "github.com/eiladin/k8s-dotenv/internal/api/v1"
	"github.com/spf13/cobra"
)

func newJobCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "job RESOURCE_NAME",
		Aliases: []string{"jobs"},
		Short:   "fetch environment configuration from job into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			list, err := v1.Jobs(opt)
			if err != nil {
				log.Fatal(err)
			}
			return list, cobra.ShellCompDirectiveDefault
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return ErrResourceNameRequired
			}
			opt.Name = args[0]
			envRes, err := v1.Job(opt)
			if err != nil {
				return err
			}

			return envRes.Write(opt)
		},
	}

	return cmd
}
