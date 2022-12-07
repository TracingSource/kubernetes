// +build !providerless

package aws

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"k8s.io/klog"
)

// Handler for aws-sdk-go that logs all requests
func awsHandlerLogger(req *request.Request) {
	service, name := awsServiceAndName(req)
	klog.V(4).Infof("AWS request: %s %s", service, name)
}

func awsSendHandlerLogger(req *request.Request) {
	service, name := awsServiceAndName(req)
	klog.V(4).Infof("AWS API Send: %s %s %v %v", service, name, req.Operation, req.Params)
}

func awsValidateResponseHandlerLogger(req *request.Request) {
	service, name := awsServiceAndName(req)
	klog.V(4).Infof("AWS API ValidateResponse: %s %s %v %v %s", service, name, req.Operation, req.Params, req.HTTPResponse.Status)
}

func awsServiceAndName(req *request.Request) (string, string) {
	service := req.ClientInfo.ServiceName

	name := "?"
	if req.Operation != nil {
		name = req.Operation.Name
	}
	return service, name
}
