package app

import (
	"fmt"

	"net/http"

	"k8s.io/kubernetes/pkg/controller/bootstrap"
)

func startBootstrapSignerController(ctx ControllerContext) (http.Handler, bool, error) {
	bsc, err := bootstrap.NewSigner(
		ctx.ClientBuilder.ClientOrDie("bootstrap-signer"),
		ctx.InformerFactory.Core().V1().Secrets(),
		ctx.InformerFactory.Core().V1().ConfigMaps(),
		bootstrap.DefaultSignerOptions(),
	)
	if err != nil {
		return nil, true, fmt.Errorf("error creating BootstrapSigner controller: %v", err)
	}
	go bsc.Run(ctx.Stop)
	return nil, true, nil
}

func startTokenCleanerController(ctx ControllerContext) (http.Handler, bool, error) {
	tcc, err := bootstrap.NewTokenCleaner(
		ctx.ClientBuilder.ClientOrDie("token-cleaner"),
		ctx.InformerFactory.Core().V1().Secrets(),
		bootstrap.DefaultTokenCleanerOptions(),
	)
	if err != nil {
		return nil, true, fmt.Errorf("error creating TokenCleaner controller: %v", err)
	}
	go tcc.Run(ctx.Stop)
	return nil, true, nil
}
