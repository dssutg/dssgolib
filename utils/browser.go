package utils

import (
	"errors"
	"os/exec"
	"runtime"
)

var ErrUnsupportedOS = errors.New("unsupported OS")

// OpenURI opens the provided URI in the OS default browser.
// Only Linux with XDG, Windows, and macOS are supported.
func OpenURI(uri string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", uri).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Start()
	case "darwin":
		return exec.Command("open", uri).Start()
	default:
		return ErrUnsupportedOS
	}
}
