package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/eiladin/k8s-dotenv/cmd/completion"
	"github.com/eiladin/k8s-dotenv/cmd/doc"
	"github.com/eiladin/k8s-dotenv/cmd/get"
	"github.com/eiladin/k8s-dotenv/pkg/client"
	"github.com/eiladin/k8s-dotenv/pkg/clioptions"
	"github.com/eiladin/k8s-dotenv/pkg/kubeclient"
	"github.com/spf13/cobra"
)

// ErrNoFilename is returned when no filename is provided.
var ErrNoFilename = errors.New("no filename provided")

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
			kubeClient, err := kubeclient.GetDefault()
			if err != nil {
				//nolint
				return err
			}

			opt.KubeClient = kubeClient

			if stdOut {
				opt.Writer = os.Stdout
			} else {
				if opt.Filename == "" {
					return ErrNoFilename
				}

				//nolint
				f, err := os.OpenFile(opt.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					return fmt.Errorf("creating output file: %w", err)
				}

				opt.Writer = f
			}

			if opt.Namespace == "" {
				if err := opt.ResolveNamespace(); err != nil {
					//nolint
					return err
				}
			}

			return nil
		},
		Version: version,
	}

	cmd.PersistentFlags().StringVarP(&opt.Namespace, "namespace", "n", "", "Namespace (default current context namespace)")
	cmd.PersistentFlags().StringVarP(&opt.Filename, "outfile", "o", ".env", "Output file")
	cmd.PersistentFlags().BoolVarP(&opt.NoExport, "no-export", "e", false, "Do not include `export` statements")
	cmd.PersistentFlags().BoolVarP(&stdOut, "console", "c", false, "Output to console")

	_ = cmd.RegisterFlagCompletionFunc("namespace",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			list, err := client.NewClient(client.WithKubeClient(opt.KubeClient)).CoreV1().NamespaceList()
			if err != nil {
				log.Fatal(err)
			}

			return list, cobra.ShellCompDirectiveDefault
		})

	cmd.AddCommand(
		completion.NewCmd(opt),
		get.NewCmd(opt),
		doc.NewCmd(opt),
	)

	root.cmd = cmd

	return root
}
