/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/eiladin/k8s-dotenv/internal/namespace"
	"github.com/spf13/cobra"
)

var namespaceName string
var contextNamespace string
var outfile string
var noExport bool

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
	}

	cmd.PersistentFlags().StringVarP(&namespaceName, "namespace", "n", "", "Namespace")
	_ = cmd.RegisterFlagCompletionFunc("namespace", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, err := namespace.GetList()
		if err != nil {
			log.Fatal(err)
		}
		return list, cobra.ShellCompDirectiveDefault
	})

	cmd.PersistentFlags().StringVarP(&outfile, "outfile", "o", ".env", "Output file")

	cmd.PersistentFlags().BoolVarP(&noExport, "no-export", "e", false, "Do not include `export` statements")

	cmd.AddCommand(
		newCompletionCmd(""),
		newGetCmd(),
	)

	root.cmd = cmd
	return root
}
