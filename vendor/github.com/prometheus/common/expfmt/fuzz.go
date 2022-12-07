// Build only when actually fuzzing
// +build gofuzz

package expfmt

import "bytes"

// Fuzz text metric parser with with github.com/dvyukov/go-fuzz:
//
//     go-fuzz-build github.com/prometheus/common/expfmt
//     go-fuzz -bin expfmt-fuzz.zip -workdir fuzz
//
// Further input samples should go in the folder fuzz/corpus.
func Fuzz(in []byte) int {
	parser := TextParser{}
	_, err := parser.TextToMetricFamilies(bytes.NewReader(in))

	if err != nil {
		return 0
	}

	return 1
}
