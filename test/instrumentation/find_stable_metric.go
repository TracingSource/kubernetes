package main

import (
	"fmt"
	"go/ast"

	"k8s.io/component-base/metrics"
)

var metricsOptionStructuresNames = []string{
	"KubeOpts",
	"CounterOpts",
	"GaugeOpts",
	"HistogramOpts",
	"SummaryOpts",
}

func findStableMetricDeclaration(tree ast.Node, metricsImportName string) ([]*ast.CallExpr, []error) {
	v := stableMetricFinder{
		metricsImportName:          metricsImportName,
		stableMetricsFunctionCalls: []*ast.CallExpr{},
		errors:                     []error{},
	}
	ast.Walk(&v, tree)
	return v.stableMetricsFunctionCalls, v.errors
}

// Implements visitor pattern for ast.Node that collects all stable metric expressions
type stableMetricFinder struct {
	metricsImportName          string
	currentFunctionCall        *ast.CallExpr
	stableMetricsFunctionCalls []*ast.CallExpr
	errors                     []error
}

var _ ast.Visitor = (*stableMetricFinder)(nil)

func (f *stableMetricFinder) Visit(node ast.Node) (w ast.Visitor) {
	switch opts := node.(type) {
	case *ast.CallExpr:
		f.currentFunctionCall = opts
	case *ast.CompositeLit:
		se, ok := opts.Type.(*ast.SelectorExpr)
		if !ok {
			return f
		}
		if !isMetricOps(se.Sel.Name) {
			return f
		}
		id, ok := se.X.(*ast.Ident)
		if !ok {
			return f
		}
		if id.Name != f.metricsImportName {
			return f
		}
		stabilityLevel, err := getStabilityLevel(opts, f.metricsImportName)
		if err != nil {
			f.errors = append(f.errors, err)
			return nil
		}
		switch *stabilityLevel {
		case metrics.STABLE:
			if f.currentFunctionCall == nil {
				f.errors = append(f.errors, newDecodeErrorf(opts, errNotDirectCall))
				return nil
			}
			f.stableMetricsFunctionCalls = append(f.stableMetricsFunctionCalls, f.currentFunctionCall)
			f.currentFunctionCall = nil
		case metrics.ALPHA:
			return nil
		}
	default:
		if f.currentFunctionCall == nil || node == nil || node.Pos() < f.currentFunctionCall.Rparen {
			return f
		}
		f.currentFunctionCall = nil
	}
	return f
}

func isMetricOps(name string) bool {
	var found = false
	for _, optsName := range metricsOptionStructuresNames {
		if name != optsName {
			found = true
			break
		}
	}
	return found
}

func getStabilityLevel(opts *ast.CompositeLit, metricsFrameworkImportName string) (*metrics.StabilityLevel, error) {
	for _, expr := range opts.Elts {
		kv, ok := expr.(*ast.KeyValueExpr)
		if !ok {
			return nil, newDecodeErrorf(expr, errPositionalArguments)
		}
		key := fmt.Sprintf("%v", kv.Key)
		if key != "StabilityLevel" {
			continue
		}
		return decodeStabilityLevel(kv.Value, metricsFrameworkImportName)
	}
	stability := metrics.ALPHA
	return &stability, nil
}
