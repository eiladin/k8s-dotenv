/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io/ioutil"
	"log"

	"github.com/eiladin/k8s-dotenv/internal/client"
	"github.com/eiladin/k8s-dotenv/internal/configmap"
	"github.com/eiladin/k8s-dotenv/internal/deployment"
	"github.com/eiladin/k8s-dotenv/internal/namespace"
	"github.com/eiladin/k8s-dotenv/internal/secret"
	"github.com/spf13/cobra"
)

var namespaceName string
var contextNamespace string
var deploymentName string
var outfile string

func Execute(version string, args []string) {
	newRootCmd(version).Execute(args)
}

type rootCmd struct {
	cmd *cobra.Command
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func newRootCmd(version string) *rootCmd {
	var root = &rootCmd{}
	var cmd = &cobra.Command{
		Use:     "k8s-dotenv",
		Short:   "Convert kubernetes secrets or configmaps to .env files",
		Long:    `k8s-dotenv takes a kubernetes secret or configmap and turns it into a .env file.`,
		Version: version,
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

	cmd.Flags().StringVarP(&namespaceName, "namespace", "n", "", "Namespace")
	_ = cmd.RegisterFlagCompletionFunc("namespace", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, err := namespace.GetList()
		if err != nil {
			log.Fatal(err)
		}
		return list, cobra.ShellCompDirectiveDefault
	})

	cmd.Flags().StringVarP(&outfile, "outfile", "o", ".env", "Output file")

	cmd.Flags().StringVarP(&deploymentName, "deployment", "d", "", "Deployment")
	_ = cmd.MarkFlagRequired("deployment")
	_ = cmd.RegisterFlagCompletionFunc("deployment", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, err := deployment.GetList(namespaceName)
		if err != nil {
			log.Fatal(err)
		}
		return list, cobra.ShellCompDirectiveDefault
	})

	cmd.AddCommand(
		newCompletionCmd(""),
	)

	root.cmd = cmd
	return root
}
