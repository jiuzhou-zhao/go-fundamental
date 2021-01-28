package pathutils

import (
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// tmpPermissionForDirectory makes the destination directory writable,
	// so that stuff can be copied recursively even if any original directory is NOT writable.
	// See https://github.com/otiai10/copy/pull/9 for more information.
	tmpPermissionForDirectory = os.FileMode(0755)
)

// IsPathExists function
func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// IsFileExists function
func IsFileExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		return !fi.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// IsDirExists function
func IsDirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		return fi.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// MakesureDirExists function
func MakesureDirExists(path string) error {
	return MustDirExists(path)
}

// MustDirExists function
func MustDirExists(path string) error {
	exists, err := IsDirExists(path)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// MakesureDirOfFileExists function
func MakesureDirOfFileExists(file string) error {
	return MakesureDirExists(filepath.Dir(file))
}

// MustDirOfFileExists function
func MustDirOfFileExists(file string) error {
	return MustDirExists(filepath.Dir(file))
}

// RemoveAll function
func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// LookPath searches for an executable named file in the
// directories named by the PATH environment variable.
// can change the PATH
func LookPath(file string) (string, error) {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = build.Default.GOPATH
	}
	newPath := strings.Join([]string{filepath.Join(goPath, "bin"), os.Getenv("PATH")},
		string(filepath.ListSeparator))

	_ = os.Setenv("PATH", newPath)

	return exec.LookPath(file)
}

// LCopy is for a symlink,
// with just creating a new symlink by replicating src symlink.
func LCopy(src, dest string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dest)
}

// DCopy is for a directory,
// with scanning contents inside the directory
// and pass everything to "copy" recursively.
func DCopy(srcdir, destdir string, info os.FileInfo) error {

	originalMode := info.Mode()

	// Make dest dir with 0755 so that everything writable.
	if err := os.MkdirAll(destdir, tmpPermissionForDirectory); err != nil {
		return err
	}
	// Recover dir mode with original one.
	defer os.Chmod(destdir, originalMode)

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())
		if err := Copy(cs, cd, content); err != nil {
			// If any error, exit immediately
			return err
		}
	}

	return nil
}

// FCopy is for just a file,
// with considering existence of parent directory
// and file permission.
func FCopy(src, dest string, info os.FileInfo) error {

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

// copy dispatches copy-funcs according to the mode.
// Because this "copy" could be called recursively,
// "info" MUST be given here, NOT nil.
func Copy(src, dest string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return LCopy(src, dest, info)
	}
	if info.IsDir() {
		return DCopy(src, dest, info)
	}
	return FCopy(src, dest, info)
}

// XCopy copies src to dest, doesn't matter if src is a directory or a file
func XCopy(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return Copy(src, dest, info)
}
