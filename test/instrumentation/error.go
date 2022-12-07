package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

const (
	errNotDirectCall        = "Opts for STABLE metric was not directly passed to new metric function"
	errPositionalArguments  = "Positional arguments are not supported"
	errStabilityLevel       = "StabilityLevel should be passed STABLE, ALPHA or removed"
	errStableSummary        = "Stable summary metric is not supported"
	errInvalidNewMetricCall = "Invalid new metric call, please ensure code compiles"
	errNonStringAttribute   = "Non string attribute it not supported"
	errFieldNotSupported    = "Field %s is not supported"
	errBuckets              = "Buckets should be set to list of floats, result from function call of prometheus.LinearBuckets or prometheus.ExponentialBuckets"
	errLabels               = "Labels were not set to list of strings"
	errImport               = `Importing using "." is not supported`
)

type decodeError struct {
	msg string
	pos token.Pos
}

func newDecodeErrorf(node ast.Node, format string, a ...interface{}) *decodeError {
	return &decodeError{
		msg: fmt.Sprintf(format, a...),
		pos: node.Pos(),
	}
}

var _ error = (*decodeError)(nil)

func (e decodeError) Error() string {
	return e.msg
}

func (e decodeError) errorWithFileInformation(fileset *token.FileSet) error {
	position := fileset.Position(e.pos)
	return fmt.Errorf("%s:%d:%d: %s", position.Filename, position.Line, position.Column, e.msg)
}
