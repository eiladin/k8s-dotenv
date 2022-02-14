package cmd

import (
	"log"

	"github.com/eiladin/k8s-dotenv/internal/client"
	"github.com/eiladin/k8s-dotenv/internal/daemonset"
	"github.com/spf13/cobra"
)

func newDaemonSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "daemonset DAEMONSET_NAME",
		Aliases: []string{"daemonsets", "ds"},
		Short:   "fetch environment configuration from daemon set into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			ns, err := client.CurrentNamespace(namespaceName)
			if err != nil {
				log.Fatal(err)
			}

			list, err := daemonset.GetList(ns)
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

			envRes, err := daemonset.Get(ns, args[0])
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
