//go:build windows

package dir

import (
	"syscall"
)

func hideFile(filename string) error {
	name, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}
	err = syscall.SetFileAttributes(name, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		return err
	}
	return nil
}
