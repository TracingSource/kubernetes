package cmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/alpha"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/upgrade"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	// Register the kubeadm configuration types because CLI flag generation
	// depends on the generated defaults.
)

// NewKubeadmCommand returns cobra.Command to run kubeadm command
func NewKubeadmCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var rootfsPath string

	cmds := &cobra.Command{
		Use:   "kubeadm",
		Short: "kubeadm: easily bootstrap a secure Kubernetes cluster",
		Long: dedent.Dedent(`

			    ┌──────────────────────────────────────────────────────────┐
			    │ KUBEADM                                                  │
			    │ Easily bootstrap a secure Kubernetes cluster             │
			    │                                                          │
			    │ Please give us feedback at:                              │
			    │ https://github.com/kubernetes/kubeadm/issues             │
			    └──────────────────────────────────────────────────────────┘

			Example usage:

			    Create a two-machine cluster with one control-plane node
			    (which controls the cluster), and one worker node
			    (where your workloads, like Pods and Deployments run).

			    ┌──────────────────────────────────────────────────────────┐
			    │ On the first machine:                                    │
			    ├──────────────────────────────────────────────────────────┤
			    │ control-plane# kubeadm init                              │
			    └──────────────────────────────────────────────────────────┘

			    ┌──────────────────────────────────────────────────────────┐
			    │ On the second machine:                                   │
			    ├──────────────────────────────────────────────────────────┤
			    │ worker# kubeadm join <arguments-returned-from-init>      │
			    └──────────────────────────────────────────────────────────┘

			    You can then repeat the second step on as many other machines as you like.

		`),
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if rootfsPath != "" {
				if err := kubeadmutil.Chroot(rootfsPath); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmds.ResetFlags()

	cmds.AddCommand(NewCmdCompletion(out, ""))
	cmds.AddCommand(NewCmdConfig(out))
	cmds.AddCommand(NewCmdInit(out, nil))
	cmds.AddCommand(NewCmdJoin(out, nil))
	cmds.AddCommand(NewCmdReset(in, out, nil))
	cmds.AddCommand(NewCmdVersion(out))
	cmds.AddCommand(NewCmdToken(out, err))
	cmds.AddCommand(upgrade.NewCmdUpgrade(out))
	cmds.AddCommand(alpha.NewCmdAlpha(in, out))

	options.AddKubeadmOtherFlags(cmds.PersistentFlags(), &rootfsPath)

	return cmds
}
