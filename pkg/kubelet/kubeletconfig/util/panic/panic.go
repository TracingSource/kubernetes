package panic

import utilruntime "k8s.io/apimachinery/pkg/util/runtime"

// HandlePanic returns a function that wraps `fn` with the utilruntime.PanicHandlers, and continues
// to bubble the panic after the PanicHandlers are called
func HandlePanic(fn func()) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				for _, fn := range utilruntime.PanicHandlers {
					fn(r)
				}
				panic(r)
			}
		}()
		// call the function
		fn()
	}
}
