package util

import (
	"os/exec"
)

// CopyDir copies the content of a folder
func CopyDir(src string, dst string) error {
	cmd := exec.Command("cp", "-r", src, dst)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
