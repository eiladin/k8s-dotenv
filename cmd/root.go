package cmd

import (
	"log"
	"os"

	"github.com/eiladin/k8s-dotenv/cmd/completion"
	"github.com/eiladin/k8s-dotenv/cmd/get"
	v1 "github.com/eiladin/k8s-dotenv/pkg/api/v1"
	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
)

var opt *options.Options = &options.Options{}
var stdOut bool = false

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
		Use:   "k8s-dotenv",
		Short: "Convert kubernetes secrets or configmaps to .env files",
		Long:  `k8s-dotenv takes a kubernetes secret or configmap and turns it into a .env file.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFlags(0)
			cs, err := client.Get()
			if err != nil {
				return err
			}
			opt.Client = cs
			if stdOut {
				opt.Writer = os.Stdout
			}

			return opt.ResolveNamespace("")
		},
		Version: version,
	}

	cmd.PersistentFlags().StringVarP(&opt.Namespace, "namespace", "n", "", "Namespace")
	_ = cmd.RegisterFlagCompletionFunc("namespace", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list, err := v1.Namespaces(opt)
		if err != nil {
			log.Fatal(err)
		}
		return list, cobra.ShellCompDirectiveDefault
	})

	cmd.PersistentFlags().StringVarP(&opt.Filename, "outfile", "o", ".env", "Output file")

	cmd.PersistentFlags().BoolVarP(&opt.NoExport, "no-export", "e", false, "Do not include `export` statements")

	cmd.PersistentFlags().BoolVarP(&stdOut, "console", "c", false, "Output to console")

	cmd.AddCommand(
		completion.NewCmd(opt),
		get.NewCmd(opt),
	)

	root.cmd = cmd
	return root
}
