package pod

import (
	v1 "github.com/eiladin/k8s-dotenv/pkg/api/v1"
	"github.com/eiladin/k8s-dotenv/pkg/errors/cmd"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
)

// NewCmd creates the `pod` command.
func NewCmd(opt *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pod RESOURCE_NAME",
		Aliases: []string{"pods", "po"},
		Short:   "fetch environment configuration from pod into a file",
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
	list, _ := v1.Pods(opt.Client, opt.Namespace)
	return list
}

func run(opt *options.Options, args []string) error {
	if len(args) == 0 {
		return cmd.ErrResourceNameRequired
	}

	res, err := v1.Pod(opt.Client, opt.Namespace, args[0])
	if err != nil {
		return err
	}

	return res.Write(opt)
}
