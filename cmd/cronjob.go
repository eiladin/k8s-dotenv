package cmd

import (
	"log"

	"github.com/eiladin/k8s-dotenv/internal/client"
	"github.com/eiladin/k8s-dotenv/internal/cronjob"
	"github.com/spf13/cobra"
)

func newCronJobCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cronjob JOB_NAME",
		Aliases: []string{"cronjobs", "cj"},
		Short:   "fetch environment configuration from cronjob into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			ns, err := client.CurrentNamespace(namespaceName)
			if err != nil {
				log.Fatal(err)
			}

			list, err := cronjob.GetList(ns)
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

			envRes, err := cronjob.Get(ns, args[0])
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
