/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/eiladin/k8s-dotenv/internal/configmap"
	"github.com/eiladin/k8s-dotenv/internal/deployment"
	"github.com/eiladin/k8s-dotenv/internal/secret"
	"github.com/spf13/cobra"
)

var namespaceName string
var deploymentName string
var outfile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-dotenv",
	Short: "Convert kubernetes secrets or configmaps to .env files",
	Long:  `k8s-dotenv takes a kubernetes secret or configmap and turns it into a .env file.`,
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

		ioutil.WriteFile(outfile, []byte(res), 0644)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&namespaceName, "namespace", "n", "default", "Namespace")
	rootCmd.Flags().StringVarP(&deploymentName, "deployment", "d", "", "Deployment")
	rootCmd.Flags().StringVarP(&outfile, "outfile", "o", ".env", "Output file")
	_ = rootCmd.MarkFlagRequired("deployment")
}
