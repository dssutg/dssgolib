//go:build !linux
// +build !linux

package logrot

import (
	"os"
)

// no-op on OSes other than Linux
func chown(name string, info os.FileInfo) error {
	return nil
}
