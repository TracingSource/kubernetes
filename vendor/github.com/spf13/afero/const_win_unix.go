// Copyright © 2016 Steve Francia <spf@spf13.com>.
//

// +build !darwin
// +build !openbsd
// +build !freebsd
// +build !dragonfly
// +build !netbsd

package afero

import (
	"syscall"
)

const BADFD = syscall.EBADFD
