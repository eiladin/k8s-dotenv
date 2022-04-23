package doc

import (
	"os"

	"github.com/eiladin/k8s-dotenv/pkg/options"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// NewCmd creates the `get` command.
func NewCmd(opt *options.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "doc",
		Short:  "generate markdown documentation",
		Hidden: true,
		RunE: func(c *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(c.Parent(), "./docs") //nolint
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			opt.Writer = os.Stdout
		},
	}

	return cmd
}
