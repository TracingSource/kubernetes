package app

import (
	"fmt"
	"io/ioutil"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	clientset "k8s.io/client-go/kubernetes"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/controller"
	serviceaccountcontroller "k8s.io/kubernetes/pkg/controller/serviceaccount"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/serviceaccount"
)

// serviceAccountTokenControllerStarter 比较特殊, ta必须先启动, 给其他controller提供权限.
//
// 这个任务无法使用常规的 client builder 对象完成, 只能手动读取sa密钥对以创建 sa token controller.
// 不能包含在常规的初始化 map 中(NewControllerInitializers()函数中可查看),
// 需要先于这个 map 中初始化函数执行(StartControllers()函数中的确是这样做的)
//
// serviceAccountTokenControllerStarter is special
// because it must run first to set up permissions for other controllers.
// It cannot use the "normal" client builder, so it tracks its own.
// It must also avoid being included in the "normal"
// init map so that it can always run first.
type serviceAccountTokenControllerStarter struct {
	rootClientBuilder controller.ControllerClientBuilder
}

// startServiceAccountTokenController ...
//
// caller: 
// 	1. cmd/kube-controller-manager/app/controllermanager.go -> StartControllers()
//  在 kcm 启动时, 选主完成后, 主进程调用.
func (c serviceAccountTokenControllerStarter) startServiceAccountTokenController(
	ctx ControllerContext,
) (http.Handler, bool, error) {
	if !ctx.IsControllerEnabled(saTokenControllerName) {
		klog.Warningf("%q is disabled", saTokenControllerName)
		return nil, false, nil
	}
	// /etc/kubernetes/pki 目录下, 只存在 sa.key 与 sa.pub 密钥对, 但是没有 crt 证书.
	if len(ctx.ComponentConfig.SAController.ServiceAccountKeyFile) == 0 {
		klog.Warningf(
			"%q is disabled because there is no private key", 
			saTokenControllerName,
		)
		return nil, false, nil
	}
	privateKey, err := keyutil.PrivateKeyFromFile(
		ctx.ComponentConfig.SAController.ServiceAccountKeyFile,
	)
	if err != nil {
		return nil, true, fmt.Errorf(
			"error reading key for service account token controller: %v", err,
		)
	}

	var rootCA []byte
	if ctx.ComponentConfig.SAController.RootCAFile != "" {
		rootCA, err = readCA(ctx.ComponentConfig.SAController.RootCAFile)
		if err != nil {
			return nil, true, fmt.Errorf(
				"error parsing root-ca-file at %s: %v", 
				ctx.ComponentConfig.SAController.RootCAFile, err,
			)
		}
	} else {
		rootCA = c.rootClientBuilder.ConfigOrDie("tokens-controller").CAData
	}

	tokenGenerator, err := serviceaccount.JWTTokenGenerator(
		serviceaccount.LegacyIssuer, privateKey,
	)
	if err != nil {
		return nil, false, fmt.Errorf("failed to build token generator: %v", err)
	}
	controller, err := serviceaccountcontroller.NewTokensController(
		ctx.InformerFactory.Core().V1().ServiceAccounts(),
		ctx.InformerFactory.Core().V1().Secrets(),
		c.rootClientBuilder.ClientOrDie("tokens-controller"),
		serviceaccountcontroller.TokensControllerOptions{
			TokenGenerator: tokenGenerator,
			RootCA:         rootCA,
		},
	)
	if err != nil {
		return nil, true, fmt.Errorf("error creating Tokens controller: %v", err)
	}
	go controller.Run(
		int(ctx.ComponentConfig.SAController.ConcurrentSATokenSyncs), ctx.Stop,
	)

	// start the first set of informers now so that other controllers can start
	ctx.InformerFactory.Start(ctx.Stop)

	return nil, true, nil
}

// readCA ioutil读取ca文件内容并返回
//
// 	@param file: /etc/kubernetes/pki/ca.crt
//
// caller:
// 	1. startServiceAccountTokenController().
func readCA(file string) ([]byte, error) {
	rootCA, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if _, err := certutil.ParseCertsPEM(rootCA); err != nil {
		return nil, err
	}

	return rootCA, err
} 

// shouldTurnOnDynamicClient 遍历 kubectl api-resources 返回的结果,
// 查找 serviceaccounts/token 资源, 确认是否拥有该资源的 create 权限. 如果有, 则返回 true.
//
// caller: 
// 	1. cmd/kube-controller-manager/app/controllermanager.go -> run()
//
func shouldTurnOnDynamicClient(client clientset.Interface) bool {
	if !utilfeature.DefaultFeatureGate.Enabled(features.TokenRequest) {
		return false
	}
	// apiResourceList 与 kubectl api-resources 结果一致.
	apiResourceList, err := client.Discovery().ServerResourcesForGroupVersion(
		corev1.SchemeGroupVersion.String(),
	)
	if err != nil {
		klog.Warningf(
			"fetch api resource lists failed, use legacy client builder: %v", 
			err,
		)
		return false
	}

	for _, resource := range apiResourceList.APIResources {
		if resource.Name == "serviceaccounts/token" &&
			resource.Group == "authentication.k8s.io" &&
			sets.NewString(resource.Verbs...).Has("create") {
			return true
		}
	}

	return false
}
