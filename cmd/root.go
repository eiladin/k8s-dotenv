package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/eiladin/k8s-dotenv/cmd/completion"
	"github.com/eiladin/k8s-dotenv/cmd/get"
	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/clioptions"
	"github.com/eiladin/k8s-dotenv/pkg/kubeclient"
	"github.com/spf13/cobra"
)

//nolint
var opt *clioptions.CLIOptions = &clioptions.CLIOptions{}

//nolint
var stdOut bool

// Execute creates the `k8s-dotenv` command with version and calls execute.
func Execute(version string, args []string) {
	newRootCmd(version).execute(args)
}

type rootCmd struct {
	cmd *cobra.Command
}

func (cmd *rootCmd) execute(args []string) {
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
			cs, err := kubeclient.Get()
			if err != nil {
				//nolint
				return err
			}
			opt.KubeClient = cs
			if stdOut {
				opt.Writer = os.Stdout
			}

			if err := opt.ResolveNamespace(""); err != nil {
				return fmt.Errorf("resolve namespace error: %w", err)
			}

			return nil
		},
		Version: version,
	}

	cmd.PersistentFlags().StringVarP(&opt.Namespace, "namespace", "n", "", "Namespace")
	_ = cmd.RegisterFlagCompletionFunc("namespace",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			list, err := client.NewClient(client.WithKubeClient(opt.KubeClient)).CoreV1().Namespaces()
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
