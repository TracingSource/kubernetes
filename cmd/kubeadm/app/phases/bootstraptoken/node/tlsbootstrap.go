package node

import (
	"fmt"

	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
)

const (
	// NodeBootstrapperClusterRoleName defines the name of the auto-bootstrapped ClusterRole for letting someone post a CSR
	// TODO: This value should be defined in an other, generic authz package instead of here
	NodeBootstrapperClusterRoleName = "system:node-bootstrapper"
	// NodeKubeletBootstrap defines the name of the ClusterRoleBinding that lets kubelets post CSRs
	NodeKubeletBootstrap = "kubeadm:kubelet-bootstrap"

	// CSRAutoApprovalClusterRoleName defines the name of the auto-bootstrapped ClusterRole for making the csrapprover controller auto-approve the CSR
	// TODO: This value should be defined in an other, generic authz package instead of here
	// Starting from v1.8, CSRAutoApprovalClusterRoleName is automatically created by the API server on startup
	CSRAutoApprovalClusterRoleName = "system:certificates.k8s.io:certificatesigningrequests:nodeclient"
	// NodeSelfCSRAutoApprovalClusterRoleName is a role defined in default 1.8 RBAC policies for automatic CSR approvals for automatically rotated node certificates
	NodeSelfCSRAutoApprovalClusterRoleName = "system:certificates.k8s.io:certificatesigningrequests:selfnodeclient"
	// NodeAutoApproveBootstrapClusterRoleBinding defines the name of the ClusterRoleBinding that makes the csrapprover approve node CSRs
	NodeAutoApproveBootstrapClusterRoleBinding = "kubeadm:node-autoapprove-bootstrap"
	// NodeAutoApproveCertificateRotationClusterRoleBinding defines name of the ClusterRoleBinding that makes the csrapprover approve node auto rotated CSRs
	NodeAutoApproveCertificateRotationClusterRoleBinding = "kubeadm:node-autoapprove-certificate-rotation"
)

// caller:
// 	1. cmd/kubeadm/app/cmd/phases/init/bootstraptoken.go -> runBootstrapToken()
// 	在 kubeadm init 过程中, 创建完成 bootstrap-token-${xxxxxx} 的 secret 资源对象后被调用.
//
// AllowBootstrapTokensToPostCSRs creates RBAC rules in a way
// that makes Node Bootstrap Tokens able to post CSRs
func AllowBootstrapTokensToPostCSRs(client clientset.Interface) error {
	fmt.Println(
		"[bootstrap-token] configured RBAC rules to allow Node Bootstrap tokens "+
		"to post CSRs in order for nodes to get long term certificate credentials",
	)

	return apiclient.CreateOrUpdateClusterRoleBinding(client, &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: NodeKubeletBootstrap,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     NodeBootstrapperClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: rbac.GroupKind,
				Name: constants.NodeBootstrapTokenAuthGroup,
			},
		},
	})
}

// caller:
// 	1. cmd/kubeadm/app/cmd/phases/init/bootstraptoken.go -> runBootstrapToken()
// 	在 kubeadm init 过程中, 创建完成 bootstrap-token-${xxxxxx} 的 secret 资源对象后被调用.
//
// AutoApproveNodeBootstrapTokens creates RBAC rules in a way
// that makes Node Bootstrap Tokens' CSR auto-approved by the csrapprover controller
func AutoApproveNodeBootstrapTokens(client clientset.Interface) error {
	fmt.Println(
		"[bootstrap-token] configured RBAC rules to allow the csrapprover controller "+
		"automatically approve CSRs from a Node Bootstrap Token",
	)

	// Always create this kubeadm-specific binding though
	return apiclient.CreateOrUpdateClusterRoleBinding(client, &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: NodeAutoApproveBootstrapClusterRoleBinding,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     CSRAutoApprovalClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: "Group",
				Name: constants.NodeBootstrapTokenAuthGroup,
			},
		},
	})
}

// caller:
// 	1. cmd/kubeadm/app/cmd/phases/init/bootstraptoken.go -> runBootstrapToken()
// 	在 kubeadm init 过程中, 创建完成 bootstrap-token-${xxxxxx} 的 secret 资源对象后被调用.
//
// AutoApproveNodeCertificateRotation creates RBAC rules in a way
// that makes Node certificate rotation CSR auto-approved by the csrapprover controller
func AutoApproveNodeCertificateRotation(client clientset.Interface) error {
	fmt.Println(
		"[bootstrap-token] configured RBAC rules to allow certificate rotation "+
		"for all node client certificates in the cluster",
	)

	return apiclient.CreateOrUpdateClusterRoleBinding(client, &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: NodeAutoApproveCertificateRotationClusterRoleBinding,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     NodeSelfCSRAutoApprovalClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: "Group",
				Name: constants.NodesGroup,
			},
		},
	})
}
