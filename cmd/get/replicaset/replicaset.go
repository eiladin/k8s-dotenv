package replicaset

import (
	v1 "github.com/eiladin/k8s-dotenv/pkg/api/v1"
	"github.com/eiladin/k8s-dotenv/pkg/errors/cmd"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
)

// NewCmd creates the `replicaset` command.
func NewCmd(opt *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "replicaset RESOURCE_NAME",
		Aliases: []string{"replicasets", "rs"},
		Short:   "fetch environment configuration from replica set into a file",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return validArgs(opt), cobra.ShellCompDirectiveDefault
		},
		RunE: func(c *cobra.Command, args []string) error {
			return run(opt, args)
		},
	}

	return cmd
}

func validArgs(opt *options.Options) []string {
	list, _ := v1.ReplicaSets(opt.Client, opt.Namespace)
	return list
}

func run(opt *options.Options, args []string) error {
	if len(args) == 0 {
		return cmd.ErrResourceNameRequired
	}

	res, err := v1.ReplicaSet(opt.Client, opt.Namespace, args[0])
	if err != nil {
		return err
	}

	return res.Write(opt)
}
