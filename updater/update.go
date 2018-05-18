package updater

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/kardianos/osext"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

// Updater is responsible for self updating a binary
type Updater struct{}

// New creates a new instance of Updater
func New() *Updater {
	return &Updater{}
}

// Apply performs an update of the current executable (or opts.TargetFile, if set) with the contents of the given io.Reader.
func (u *Updater) Apply(reader io.Reader) error {
	tmpPath := os.TempDir()
	err := archiver.TarGz.Read(reader, tmpPath)
	if err != nil {
		return errors.Wrap(err, "Could not unzip the asset")
	}

	targetPath, err := osext.Executable()
	if err != nil {
		return errors.Wrap(err, "Could not get the executable path")
	}

	updateDir := filepath.Dir(targetPath)
	filename := filepath.Base(targetPath)

	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return errors.Wrap(err, "Could not open the file")
	}

	defer fp.Close()

	newBytes, err := ioutil.ReadFile(path.Join(tmpPath, filename))
	if err != nil {
		return errors.Wrap(err, "Could not find new binary")
	}

	_, err = io.Copy(fp, bytes.NewReader(newBytes))

	// if we don't call fp.Close(), windows won't let us move the new executable
	// because the file will still be "in use"
	fp.Close()

	oldPath := filepath.Join(updateDir, fmt.Sprintf(".%s.old", filename))

	// delete any existing old exec file - this is necessary on Windows for two reasons:
	// 1. after a successful update, Windows can't remove the .old file because the process is still running
	// 2. windows rename operations fail if the destination file already exists
	_ = os.Remove(oldPath)

	// move the existing executable to a new file in the same directory
	err = os.Rename(targetPath, oldPath)
	if err != nil {
		return errors.Wrap(err, "Could move the new binary to path")
	}

	// move the new exectuable in to become the new program
	err = os.Rename(newPath, targetPath)
	if err != nil {
		// move unsuccessful
		//
		// The filesystem is now in a bad state. We have successfully
		// moved the existing binary to a new location, but we couldn't move the new
		// binary to take its place. That means there is no file where the current executable binary
		// used to be!
		// Try to rollback by restoring the old binary to its original path.
		err = os.Rename(oldPath, targetPath)
		if err != nil {
			return errors.Wrap(err, "There was an error while renaming the binary")
		}
	}

	errRemove := os.Remove(oldPath)
	// windows has trouble with removing old binaries, so hide it instead
	if errRemove != nil {
		_ = hideFile(oldPath)
	}

	return nil
}
