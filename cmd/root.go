/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io/ioutil"
	"log"

	"github.com/eiladin/k8s-dotenv/internal/configmap"
	"github.com/eiladin/k8s-dotenv/internal/deployment"
	"github.com/eiladin/k8s-dotenv/internal/secret"
	"github.com/spf13/cobra"
)

var namespaceName string
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

	cmd.Flags().StringVarP(&namespaceName, "namespace", "n", "default", "Namespace")
	cmd.Flags().StringVarP(&deploymentName, "deployment", "d", "", "Deployment")
	cmd.Flags().StringVarP(&outfile, "outfile", "o", ".env", "Output file")
	_ = cmd.MarkFlagRequired("deployment")

	root.cmd = cmd
	return root
}
