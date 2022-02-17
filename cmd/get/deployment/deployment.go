package deployment

import (
	"log"
	"os"

	v1 "github.com/eiladin/k8s-dotenv/internal/api/v1"
	"github.com/eiladin/k8s-dotenv/internal/errors/cmd"
	"github.com/eiladin/k8s-dotenv/internal/options"
	"github.com/spf13/cobra"
)

func NewCmd(opt *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deployment RESOURCE_NAME",
		Aliases: []string{"deployments", "deploy"},
		Short:   "fetch environment configuration from deployment into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			list, err := v1.Deployments(opt)

			if err != nil {
				log.Fatal(err)
			}
			return list, cobra.ShellCompDirectiveDefault
		},
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.ErrResourceNameRequired
			}
			opt.Name = args[0]
			res, err := v1.Deployment(opt)
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
