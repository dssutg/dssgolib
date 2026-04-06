// Package sh provides convenient shell- and command-line-oriented utilities for
// Go programs. It wraps common file operations and subprocess invocations such
// as running commands, piping output, changing directories, copying files, and
// creating ZIP archives with simple, panic-on-error semantics to streamline
// scripting tasks.
package sh

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunCommandWithEnv executes the named command with the given arguments and
// environment variables. Standard input, output, and error are inherited from
// the parent process. On failure, the function prints an error to stderr and
// exits the program with status 1.
func RunCommandWithEnv(command string, args []string, env []string) {
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Command failed: %s, Error: %s", command, err)
		os.Exit(1)
	}
}

// RunCommand executes the named command with the given arguments. On failure,
// it prints an error and exits with status 1.
func RunCommand(command string, args ...string) {
	RunCommandWithEnv(command, args, []string{})
}

// RunSimpleCommand splits the input string on whitespace and executes the
// resulting command and arguments. If the string is empty, it does nothing.
func RunSimpleCommand(command string) {
	parts := strings.Split(command, " ")

	if len(parts) == 0 {
		return
	}

	RunCommand(parts[0], parts[1:]...)
}

// PipeCommand runs the command with arguments, captures its combined stdout and
// stderr, and returns the output as a string. On error, it prints the error and
// exits with status 1.
func PipeCommand(command string, args ...string) string {
	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Command failed: %s, Error: %s", command, err)
		os.Exit(1)
	}

	return string(output)
}

// PipeSimpleCommand splits the input on whitespace, runs the resulting command,
// and returns its combined output.
func PipeSimpleCommand(command string) string {
	parts := strings.Split(command, " ")

	if len(parts) == 0 {
		return ""
	}

	return PipeCommand(parts[0], parts[1:]...)
}

// GetCommandOutput runs the external command named by `name` with the given
// arguments and returns its standard output, standard error, and any execution
// error.
//
// Parameters:
// - name: command executable name or path.
// - arg: variadic list of arguments passed to the command.
//
// Returns:
//   - out: captured standard output as a string.
//   - errOut: captured standard error as a string.
//   - err: error returned from cmd.Run(); non-nil when the command fails to start
//     or exits with a non-zero status. Note that a non-nil err does not prevent
//     out and errOut from containing partial command output.
//
// Behavior notes:
// - The command's stdin is connected to the current process's stdin.
// - stdout and stderr are captured into separate buffers and returned as strings.
// - The function does not perform any trimming on the returned output strings.
func GetCommandOutput(name string, arg ...string) (out string, errOut string, err error) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	cmd := exec.Command(name, arg...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()

	out = outBuf.String()
	errOut = errBuf.String()

	return
}

// Cd changes the current working directory to dir. On error, it prints the
// error and exits with status 1.
func Cd(dir string) {
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// Rmrf removes each specified path (file or directory) recursively On error,.
// it panics via FatalOnErr .
func Rmrf(paths ...string) {
	for _, path := range paths {
		FatalOnErr(os.RemoveAll(path))
	}
}

// RemoveDirEntries walks the directory with the provided name and removes each entry
// if shouldRemove returns true for it. Any first error is returned.
func RemoveDirEntries(name string, shouldRemove func(entry os.DirEntry) bool) error {
	entries, err := os.ReadDir(name)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if shouldRemove(entry) {
			if err := os.RemoveAll(entry.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

// Mkdir creates directories (and parents) for each specified path with
// permission 0755. Panics on error.
func Mkdir(paths ...string) {
	for _, path := range paths {
		FatalOnErr(os.MkdirAll(path, 0o755))
	}
}

// Echo prints its arguments to stdout, separated by spaces, with a newline.
func Echo(parts ...any) {
	fmt.Println(parts...)
}

// Cat reads the contents of each file in order, concatenates them, and returns
// the combined string. Panics if any read operation fails.
func Cat(paths ...string) string {
	var builder strings.Builder

	for _, path := range paths {
		content, err := os.ReadFile(path)
		FatalOnErr(err)

		builder.Write(content)
	}

	return builder.String()
}

// FatalOnErr checks err and, if non-nil, prints it to stderr and exits the
// program with status 1.
func FatalOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// CopyFile copies a file or directory from src to dst.
// If src is a directory, it recursively copies all contents.
// If dst is a directory, the source name is appended to dst.
// Returns an error on failure.
func CopyFile(src, dst string) error {
	// Get the file info for the source.
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Check if the destination is a directory.
	var dstIsDir bool
	dstInfo, err := os.Stat(dst)
	switch {
	case err == nil:
		dstIsDir = dstInfo.IsDir()
	case os.IsNotExist(err):
		// If the destination does not exist, it will be created later if needed.
		dstIsDir = true
	default:
		return err
	}

	// If the source is a directory.
	if srcInfo.IsDir() {
		if !dstIsDir {
			return os.ErrInvalid // Cannot copy a directory to a file
		}

		// Create the destination directory if it doesn't exist.
		err = os.MkdirAll(dst, os.ModePerm)
		if err != nil {
			return err
		}

		// Read the contents of the source directory.
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}

		// Iterate over the entries in the directory.
		for _, entry := range entries {
			srcPath := filepath.Join(src, entry.Name())
			dstPath := filepath.Join(dst, entry.Name())

			// Recursively copy the contents.
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	} else {
		// If the source is a file.
		if dstIsDir {
			// If the destination is a directory, create the destination file path.
			dst = filepath.Join(dst, srcInfo.Name())
		}

		// Copy the file.
		sourceFile, err := os.Open(src)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		destinationFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// Cp is a shorthand for CopyFile that exits on error.
func Cp(src, dst string) {
	FatalOnErr(CopyFile(src, dst))
}

// CreateZipArchive creates a ZIP archive at target containing the file or
// directory source. Directories are added recursively, preserving file
// permissions and modification times. Returns an error on failure.
func CreateZipArchive(target, source string) error {
	// Create the zip file.
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Get absolute path for the source so we can compute relative paths later.
	absSource, err := filepath.Abs(source)
	if err != nil {
		return err
	}

	// Check if source exists.
	sourceInfo, err := os.Stat(absSource)
	if err != nil {
		return err
	}

	// If source is a directory, we want to include its contents recursively.
	// If it's a file, we add that file.
	if sourceInfo.IsDir() {
		err = filepath.Walk(absSource, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Compute the relative path to the source folder.
			relPath, err := filepath.Rel(filepath.Dir(absSource), path)
			if err != nil {
				return err
			}
			// For directories, we need to end with "/" so that unzip interprets them correctly.
			if info.IsDir() {
				if relPath == "." || relPath == ".." {
					// Skip adding
					return nil
				}
				relPath += "/"
			}
			return addFileOrDirToZip(zipWriter, path, relPath, info)
		})
		if err != nil {
			return err
		}
	} else {
		// For a single file, set the relative name to just the base of the file.
		relName := filepath.Base(absSource)
		if err := addFileOrDirToZip(zipWriter, absSource, relName, sourceInfo); err != nil {
			return err
		}
	}

	return nil
}

// addFileOrDirToZip is a helper that adds a single file or directory entry to
// the given zip.Writer. For files, it writes the file's contents; for
// directories, it writes an empty directory entry. It preserves permissions and
// timestamps from the FileInfo.
func addFileOrDirToZip(zipWriter *zip.Writer, filePath, relPath string, info os.FileInfo) error {
	// Create header based on the FileInfo structure.
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Set header name to the relative path.
	header.Name = relPath

	// Preserve file permissions by using the OS-specific mode bits.
	header.Method = zip.Deflate // Use Deflate for maximum compression

	// Set the modified time.
	header.Modified = info.ModTime()

	// If info is a directory, no content is written.
	if info.IsDir() {
		// The header for a directory ends with '/'.
		// We have already ensured that by appending "/" in the caller.
		_, err = zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		return nil
	}

	// For files, create the writer and open the file to write its contents.
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the file content into the zip entry.
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}

func Home() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain user home directory path: %v", err)
		os.Exit(1)
	}
	return home
}
