package goupdater

import (
	"io"

	"github.com/italolelis/goupdater/updater"
	"github.com/pkg/errors"
)

// Resolver represents an upstream that will be called when checking for a new version
type Resolver interface {
	Update(currentVersion string) (io.ReadCloser, error)
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error                     { return nil }
func (nopCloser) Read(p []byte) (n int, err error) { return 0, nil }

// Update updates the current binary with the chosen resolver
func Update(resolver Resolver, currentVersion string) (bool, error) {
	reader, err := resolver.Update(currentVersion)
	if err != nil {
		return false, nil
	}
	defer reader.Close()

	githubUpdater := updater.New()
	err = githubUpdater.Apply(reader)
	if err != nil {
		return false, errors.Wrap(err, "could not apply the update")
	}

	return true, nil
}
