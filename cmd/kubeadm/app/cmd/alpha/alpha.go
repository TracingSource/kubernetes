package alpha

import (
	"io"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubernetes/cmd/kubeadm/app/cmd/util"
)

// NewCmdAlpha returns "kubeadm alpha" command.
func NewCmdAlpha(in io.Reader, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alpha",
		Short: "Kubeadm experimental sub-commands",
	}

	cmd.AddCommand(newCmdCertsUtility(out))
	cmd.AddCommand(newCmdKubeletUtility())
	cmd.AddCommand(newCmdKubeConfigUtility(out))
	cmd.AddCommand(NewCmdSelfhosting(in))

	// TODO: This command should be removed as soon as the kubeadm init phase refactoring is completed.
	//		 current phases implemented as cobra.Commands should become workflow.Phases, while other utilities
	// 		 hosted under kubeadm alpha phases command should found a new home under kubeadm alpha (without phases)
	cmd.AddCommand(newCmdPhase(out))

	return cmd
}

func newCmdPhase(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phase",
		Short: "Invoke subsets of kubeadm functions separately for a manual install",
		Long:  cmdutil.MacroCommandLongDescription,
	}

	return cmd
}
