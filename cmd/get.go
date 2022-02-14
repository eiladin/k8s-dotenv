package cmd

import (
	"io/ioutil"
	"log"

	"github.com/eiladin/k8s-dotenv/internal/client"
	"github.com/eiladin/k8s-dotenv/internal/configmap"
	"github.com/eiladin/k8s-dotenv/internal/deployment"
	"github.com/eiladin/k8s-dotenv/internal/secret"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get DEPLOYMENT_NAME",
		Short: "fetch secrets and configmaps into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			ns := namespaceName
			var err error
			if ns == "" {
				ns, err = client.CurrentNamespace()
				if err != nil {
					log.Fatal(err)
				}
			}

			list, err := deployment.GetList(ns)
			if err != nil {
				log.Fatal(err)
			}
			return list, cobra.ShellCompDirectiveDefault
		},
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			contextNamespace, err = client.CurrentNamespace()
			if err != nil && namespaceName == "" {
				log.Fatal(err)
			}

			if namespaceName == "" {
				namespaceName = contextNamespace
			}

			res := ""
			secrets, configmaps, err := deployment.Get(namespaceName, deploymentName)
			if err != nil {
				log.Fatal(err)
			}

			for _, s := range secrets {
				secretVal, err := secret.Get(namespaceName, s)
				if err != nil {
					log.Fatal(err)
				}
				res += secretVal
			}

			for _, c := range configmaps {
				configmapVal, err := configmap.Get(namespaceName, c)
				if err != nil {
					log.Fatal(err)
				}
				res += configmapVal
			}

			err = ioutil.WriteFile(outfile, []byte(res), 0644)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
