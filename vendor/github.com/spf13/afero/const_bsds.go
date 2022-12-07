// Copyright © 2016 Steve Francia <spf@spf13.com>.
//


// +build darwin openbsd freebsd netbsd dragonfly

package afero

import (
	"syscall"
)

const BADFD = syscall.EBADF
