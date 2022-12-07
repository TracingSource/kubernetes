/*
 *
 * Copyright 2018 gRPC authors.
 *

 *
 */

package binarylog

import (
	"errors"
	"strings"
)

// parseMethodName splits service and method from the input. It expects format
// "/service/method".
//
// TODO: move to internal/grpcutil.
func parseMethodName(methodName string) (service, method string, _ error) {
	if !strings.HasPrefix(methodName, "/") {
		return "", "", errors.New("invalid method name: should start with /")
	}
	methodName = methodName[1:]

	pos := strings.LastIndex(methodName, "/")
	if pos < 0 {
		return "", "", errors.New("invalid method name: suffix /method is missing")
	}
	return methodName[:pos], methodName[pos+1:], nil
}
