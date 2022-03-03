package daemonset

import (
	"errors"
	"fmt"

	v1 "github.com/eiladin/k8s-dotenv/pkg/api/v1"
	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
)

// ErrResourceNameRequired is returned when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")

func newRunError(err error) error {
	return fmt.Errorf("daemonset error: %w", err)
}

// NewCmd creates the `daemonset` command.
func NewCmd(opt *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "daemonset RESOURCE_NAME",
		Aliases: []string{"daemonsets", "ds"},
		Short:   "fetch environment configuration from daemon set into a file",
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
	list, _ := v1.DaemonSets(opt.Client, opt.Namespace)

	return list
}

func run(opt *options.Options, args []string) error {
	if len(args) == 0 {
		return ErrResourceNameRequired
	}

	res, err := v1.DaemonSet(opt.Client, opt.Namespace, args[0])
	if err != nil {
		return newRunError(err)
	}

	if err := res.Write(opt); err != nil {
		return newRunError(err)
	}

	return nil
}
