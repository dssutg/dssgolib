package utils

import (
	"encoding/gob"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const (
	FilePerm = 0o644
	DirPerm  = 0o755
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func ReadJSONFile(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(v)
}

func WriteJSONFile(path string, v any, perm os.FileMode) error {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, perm)
}

func ReadGOBFile(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return gob.NewDecoder(f).Decode(v)
}

func WriteGOBFile(path string, v any, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	return gob.NewEncoder(f).Encode(v)
}

// RemoveFilenameExtension returns filename without extension.
func RemoveFilenameExtension(filename string) string {
	extension := filepath.Ext(filename)
	return strings.TrimSuffix(filename, extension)
}

// IsFilenameWindowsReserved reports whether the provided filename
// is a reserved name on Windows. Such a name is not allowed to be used
// as filename. The provided name must be a trimmed base name only (no path),
// no trailing spaces/dots.
//
// The check is case-insensitive for these names:
//   - CON, PRN, AUX, NUL
//   - COM1, COM2, COM3, COM4, COM5, COM6, COM7, COM8, COM9
//   - LPT1, LPT2, LPT3, LPT4, LPT5, LPT6, LPT7, LPT8, LPT9
func IsFilenameWindowsReserved(name string) bool {
	// Match the filename as fast as possible. This is not pretty but fast.
	// It can be even faster though if we utilized SIMD and bit tricks.
	switch len(name) {
	case 3:
		switch name[0] {
		case 'c', 'C': // CON
			if (name[1] == 'o' || name[1] == 'O') && (name[2] == 'n' || name[2] == 'N') {
				return true
			}

		case 'p', 'P': // PRN
			if (name[1] == 'r' || name[1] == 'R') && (name[2] == 'n' || name[2] == 'N') {
				return true
			}

		case 'a', 'A': // AUX
			if (name[1] == 'u' || name[1] == 'U') && (name[2] == 'x' || name[2] == 'X') {
				return true
			}

		case 'n', 'N': // NUL
			if (name[1] == 'u' || name[1] == 'U') && (name[2] == 'l' || name[2] == 'L') {
				return true
			}
		}

	case 4:
		switch name[0] {
		case 'c', 'C': // COM1, COM2, COM3, COM4, COM5, COM6, COM7, COM8, COM9
			if (name[1] == 'o' || name[1] == 'O') && (name[2] == 'm' || name[2] == 'M') {
				num := name[3]
				if num >= '1' && num <= '9' {
					return true
				}
			}

		case 'l', 'L': // LPT1, LPT2, LPT3, LPT4, LPT5, LPT6, LPT7, LPT8, LPT9
			if (name[1] == 'p' || name[1] == 'P') && (name[2] == 't' || name[2] == 'T') {
				num := name[3]
				if num >= '1' && num <= '9' {
					return true
				}
			}
		}
	}

	return false
}

// IsIllegalWindowsFilenameChar reports whether the provided
// character cannot be part of filename on Windows.
//
// The invalid characters are:
// '/', '\\', '<', '>', ':', '"', '|', '?', '*'.
func IsIllegalWindowsFilenameChar[T RuneOrByte](c T) bool {
	switch c {
	case '/', '\\', '<', '>', ':', '"', '|', '?', '*':
		return true
	default:
		return false
	}
}

// SanitizeStringToFilename converts s to filename with appended extension.
// The final name is at most 255 bytes. If the sanitied name is empty,
// "file" is returned. The extension must be a valid non-empty ASCII string
// that starts with dot.
func SanitizeStringToFilename(s, extension string) string {
	var b strings.Builder

	b.Grow(len(s))

	for _, r := range s {
		switch {
		case IsASCIIControl(r):
			continue
		case IsIllegalWindowsFilenameChar(r):
			b.WriteByte(' ')
		default:
			b.WriteRune(r)
		}
	}

	base := b.String()
	base = StringCollapseSpaces(base)
	base = TruncateStringBytes(base, 255-len(extension))
	if base == "" {
		base = "file"
	}

	base += extension

	return base
}
