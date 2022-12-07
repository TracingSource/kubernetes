package generators

import (
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// extractBoolTagOrDie gets the comment-tags for the key and asserts that, if
// it exists, the value is boolean.  If the tag did not exist, it returns
// false.
func extractBoolTagOrDie(key string, lines []string) bool {
	val, err := types.ExtractSingleBoolCommentTag("+", key, false, lines)
	if err != nil {
		klog.Fatalf(err.Error())
	}
	return val
}
