// +build !cgo

package util

func getRunningSlice() (string, error) { return "", ErrNoCGO }

func runningFromSystemService() (bool, error) { return false, ErrNoCGO }

func currentUnitName() (string, error) { return "", ErrNoCGO }
