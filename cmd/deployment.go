package cmd

import (
	"log"

	"github.com/eiladin/k8s-dotenv/internal/client"
	"github.com/eiladin/k8s-dotenv/internal/deployment"
	"github.com/spf13/cobra"
)

func newDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deployment DEPLOYMENT_NAME",
		Aliases: []string{"deployments", "deploy"},
		Short:   "fetch environment configuration from deployment into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			ns, err := client.CurrentNamespace(namespaceName)
			if err != nil {
				log.Fatal(err)
			}

			list, err := deployment.GetList(ns)
			if err != nil {
				log.Fatal(err)
			}
			return list, cobra.ShellCompDirectiveDefault
		},
		Run: func(cmd *cobra.Command, args []string) {
			ns, err := client.CurrentNamespace(namespaceName)
			if err != nil {
				log.Fatal(err)
			}

			envRes, err := deployment.Get(ns, args[0])
			if err != nil {
				log.Fatal(err)
			}

			err = envRes.Write(ns, !noExport, outfile)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
