package auth

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// NewCmdAuth returns an initialized Command instance for 'auth' sub command
func NewCmdAuth(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "auth",
		Short: "Inspect authorization",
		Long:  `Inspect authorization`,
		Run:   cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}

	cmds.AddCommand(NewCmdCanI(f, streams))
	cmds.AddCommand(NewCmdReconcile(f, streams))

	return cmds
}
