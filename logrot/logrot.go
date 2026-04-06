// Package logrot implements log rotation.
package logrot

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

const (
	backupTimeFormat = "2006-01-02T15-04-05.000"
	compressSuffix   = ".gz"
	defaultMaxSize   = 100 * 1024 * 1024 // 100 MB
)

// Rotator is an io.WriteCloser that writes to the specified filename.
//
// Rotator opens or creates the logfile on first Write.  If the file exists and
// is less than MaxSize bytes, rotor will open and append to that file.
// If the file exists and its size is >= MaxSize bytes, the file is renamed
// by putting the current time in a timestamp in the name immediately before the
// file's extension (or the end of the filename if there's no extension). A new
// log file is then created using original filename.
//
// Whenever a write would cause the current log file exceed MaxSize bytes,
// the current file is closed, renamed, and a new log file created with the
// original name. Thus, the filename you give Rotator is always the "current" log
// file.
//
// Backups use the log file name given to Rotator, in the form
// `name-timestamp.ext` where name is the filename without the extension,
// timestamp is the time at which the log was rotated formatted with the
// time.Time format of `2006-01-02T15-04-05.000` and the extension is the
// original extension.  For example, if your Rotator.Filename is
// `/var/log/foo/server.log`, a backup created at 6:30pm on Nov 11 2016 would
// use the filename `/var/log/foo/server-2016-11-04T18-30-00.000.log`
//
// # Cleaning Up Old Log Files
//
// Whenever a new logfile gets created, old log files may be deleted.  The most
// recent files according to the encoded timestamp will be retained, up to a
// number equal to MaxBackups (or all of them if MaxBackups is 0).  Any files
// with an encoded timestamp older than MaxAge days are deleted, regardless of
// MaxBackups.  Note that the time encoded in the timestamp is the rotation
// time, which may differ from the last time that file was written to.
//
// If MaxBackups and MaxAge are both 0, no old log files will be deleted.
type Rotator struct {
	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-rotor.log in
	// os.TempDir() if empty.
	Filename string

	// MaxSize is the maximum size in bytes of the log file before it gets
	// rotated. It defaults to [defaultMaxSize] bytes.
	MaxSize int

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool

	size int64
	file *os.File
	sync.Mutex

	millCh    chan bool
	startMill sync.Once
}

// logInfo is a convenience struct to return the filename and its embedded
// timestamp.
type logInfo struct {
	timestamp time.Time
	os.DirEntry
}

// Write implements io.Writer.  If a write would cause the log file to be larger
// than MaxSize, the file is closed, renamed to include a timestamp of the
// current time, and a new log file is created using the original log file name.
// If the length of the write is greater than MaxSize, an error is returned.
func (r *Rotator) Write(p []byte) (int, error) {
	r.Lock()
	defer r.Unlock()

	writeLen := int64(len(p))
	if writeLen > r.max() {
		return 0, fmt.Errorf(
			"write length %d exceeds maximum file size %d", writeLen, r.max(),
		)
	}

	if r.file == nil {
		writeLen := len(p)

		r.mill()

		filename := r.filename()
		info, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return 0, r.openNew()
		}
		if err != nil {
			return 0, fmt.Errorf("error getting log file info: %w", err)
		}

		if info.Size()+int64(writeLen) >= r.max() {
			return 0, r.rotate()
		}

		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			// if we fail to open the old log file for some reason, just ignore
			// it and open a new log file.
			return 0, r.openNew()
		}
		r.file = file
		r.size = info.Size()
	}

	if r.size+writeLen > r.max() {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := r.file.Write(p)
	r.size += int64(n)

	return n, err
}

// Close implements io.Closer, and closes the current logfile.
func (r *Rotator) Close() error {
	r.Lock()
	defer r.Unlock()
	return r.close()
}

// close closes the file if it is open.
func (r *Rotator) close() error {
	if r.file == nil {
		return nil
	}
	err := r.file.Close()
	r.file = nil
	return err
}

// Rotate causes Rotor to close the existing log file and immediately create a
// new one.  This is a helper function for applications that want to initiate
// rotations outside of the normal rotation rules, such as in response to
// SIGHUP.  After rotating, this initiates compression and removal of old log
// files according to the configuration.
func (r *Rotator) Rotate() error {
	r.Lock()
	defer r.Unlock()
	return r.rotate()
}

// rotate closes the current file, moves it aside with a timestamp in the name,
// (if it exists), opens a new file with the original filename, and then runs
// post-rotation processing and removal.
func (r *Rotator) rotate() error {
	if err := r.close(); err != nil {
		return err
	}
	if err := r.openNew(); err != nil {
		return err
	}
	r.mill()
	return nil
}

// openNew opens a new log file for writing, moving any old log file out of the
// way.  This methods assumes the file has already been closed.
func (r *Rotator) openNew() error {
	err := os.MkdirAll(r.dir(), 0o755)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %w", err)
	}

	name := r.filename()
	mode := os.FileMode(0o600)
	info, err := os.Stat(name)
	if err == nil {
		// Copy the mode off the old logfile.
		mode = info.Mode()

		// Create a new filename from the given name, inserting a timestamp
		// between the filename and the extension, using the local time if requested
		// (otherwise UTC).
		dir := filepath.Dir(name)
		filename := filepath.Base(name)
		ext := filepath.Ext(filename)
		prefix := filename[:len(filename)-len(ext)]
		t := time.Now()
		if !r.LocalTime {
			t = t.UTC()
		}
		timestamp := t.Format(backupTimeFormat)
		newname := filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))

		err = os.Rename(name, newname)
		if err != nil {
			return fmt.Errorf("can't rename log file: %w", err)
		}

		// this is a no-op anywhere but linux
		err = chown(name, info)
		if err != nil {
			return err
		}
	}

	// we use truncate here because this should only get called when we've moved
	// the file ourselves. if someone else creates the file in the meantime,
	// just wipe out the contents.
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %w", err)
	}
	r.file = f
	r.size = 0
	return nil
}

// filename generates the name of the logfile from the current time.
func (r *Rotator) filename() string {
	if r.Filename != "" {
		return r.Filename
	}
	name := filepath.Base(os.Args[0]) + "-rotor.log"
	return filepath.Join(os.TempDir(), name)
}

// mill performs post-rotation compression and removal of stale log files,
// starting the mill goroutine if necessary.
func (r *Rotator) mill() {
	r.startMill.Do(func() {
		r.millCh = make(chan bool, 1)

		// Run in a goroutine to manage post-rotation compression and removal
		// of old log files.
		go func() {
			for range r.millCh {
				// Perform compression and removal of stale log files.
				// Log files are compressed if enabled via configuration and old log
				// files are removed, keeping at most l.MaxBackups files, as long as
				// none of them are older than MaxAge.
				if r.MaxBackups == 0 && r.MaxAge == 0 && !r.Compress {
					continue
				}

				// Get the list of backup log files stored in the same
				// directory as the current log file, sorted by ModTime
				files, filesErr := os.ReadDir(r.dir())
				if filesErr != nil {
					continue
				}
				logFiles := []logInfo{}

				filename := filepath.Base(r.filename())
				ext := filepath.Ext(filename)
				prefix := filename[:len(filename)-len(ext)] + "-"

				for _, f := range files {
					if f.IsDir() {
						continue
					}
					t, err := r.timeFromName(f.Name(), prefix, ext)
					if err == nil {
						logFiles = append(logFiles, logInfo{t, f})
						continue
					}
					t, err = r.timeFromName(f.Name(), prefix, ext+compressSuffix)
					if err == nil {
						logFiles = append(logFiles, logInfo{t, f})
						continue
					}
					// error parsing means that the suffix at the end was not generated
					// by rotor, and therefore it's not a backup file.
				}
				slices.SortFunc(logFiles, func(a, b logInfo) int {
					return b.timestamp.Compare(a.timestamp)
				})

				var compress, remove []logInfo

				if r.MaxBackups > 0 && r.MaxBackups < len(logFiles) {
					preserved := make(map[string]bool)
					var remaining []logInfo
					for _, f := range logFiles {
						fn := strings.TrimSuffix(f.Name(), compressSuffix)
						preserved[fn] = true
						if len(preserved) > r.MaxBackups {
							remove = append(remove, f)
						} else {
							remaining = append(remaining, f)
						}
					}
					logFiles = remaining
				}

				if r.MaxAge > 0 {
					diff := time.Duration(int64(24*time.Hour) * int64(r.MaxAge))
					cutoff := time.Now().Add(-1 * diff)
					var remaining []logInfo
					for _, f := range logFiles {
						if f.timestamp.Before(cutoff) {
							remove = append(remove, f)
						} else {
							remaining = append(remaining, f)
						}
					}
					logFiles = remaining
				}

				if r.Compress {
					for _, f := range logFiles {
						if !strings.HasSuffix(f.Name(), compressSuffix) {
							compress = append(compress, f)
						}
					}
				}

				for _, f := range remove {
					errRemove := os.Remove(filepath.Join(r.dir(), f.Name()))
					if filesErr == nil && errRemove != nil {
						filesErr = errRemove
					}
				}

				for _, f := range compress {
					fn := filepath.Join(r.dir(), f.Name())
					errCompress := func(dst string) error {
						var err error
						f, err := os.Open(fn)
						if err != nil {
							return fmt.Errorf("failed to open log file: %w", err)
						}
						defer f.Close()
						fi, err := os.Stat(fn)
						if err != nil {
							return fmt.Errorf("failed to stat log file: %w", err)
						}
						err = chown(dst, fi)
						if err != nil {
							return fmt.Errorf("failed to chown compressed log file: %w", err)
						}
						gzf, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
						if err != nil {
							return fmt.Errorf("failed to open compressed log file: %w", err)
						}
						defer gzf.Close()
						gz := gzip.NewWriter(gzf)
						defer func() {
							if err != nil {
								os.Remove(dst)
								err = fmt.Errorf("failed to compress log file: %w", err)
							}
						}()
						if _, err := io.Copy(gz, f); err != nil {
							return err
						}
						if err := gz.Close(); err != nil {
							return err
						}
						if err := gzf.Close(); err != nil {
							return err
						}
						if err := f.Close(); err != nil {
							return err
						}
						if err := os.Remove(fn); err != nil {
							return err
						}
						return nil
					}(fn + compressSuffix)
					if filesErr == nil && errCompress != nil {
						filesErr = errCompress
					}
				}
			}
		}()
	})

	select {
	case r.millCh <- true:
	default:
	}
}

// timeFromName extracts the formatted time from the filename by stripping off
// the filename's prefix and extension. This prevents someone's filename from
// confusing time.parse.
func (r *Rotator) timeFromName(filename, prefix, ext string) (time.Time, error) {
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, errors.New("mismatched prefix")
	}
	if !strings.HasSuffix(filename, ext) {
		return time.Time{}, errors.New("mismatched extension")
	}
	ts := filename[len(prefix) : len(filename)-len(ext)]
	return time.Parse(backupTimeFormat, ts)
}

// max returns the maximum size in bytes of log files before rolling.
func (r *Rotator) max() int64 {
	if r.MaxSize == 0 {
		return defaultMaxSize
	}
	return int64(r.MaxSize)
}

// dir returns the directory for the current filename.
func (r *Rotator) dir() string {
	return filepath.Dir(r.filename())
}
