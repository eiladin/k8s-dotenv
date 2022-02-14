package cmd

import (
	"log"

	"github.com/eiladin/k8s-dotenv/internal/client"
	"github.com/eiladin/k8s-dotenv/internal/job"
	"github.com/spf13/cobra"
)

func newJobCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "job JOB_NAME",
		Aliases: []string{"jobs"},
		Short:   "fetch environment configuration from job into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			ns, err := client.CurrentNamespace(namespaceName)
			if err != nil {
				log.Fatal(err)
			}

			list, err := job.GetList(ns)
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

			envRes, err := job.Get(ns, args[0])
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
