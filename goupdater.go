// Package goupdater helps you to easily build self-updating programs
//
// Goupdater makes it easier for you to update your Go programs (or other single-file targets).
// A program can update itself by replacing its executable file with a new version.
// It provides the flexibility to implement different updating user experiences like auto-updating,
// or manual user-initiated updates.
package goupdater

import (
	"context"
	"io"
	"strings"

	"github.com/italolelis/goupdater/updater"
	"github.com/pkg/errors"
)

// Resolver represents an upstream that will be called when checking for a new version
type Resolver interface {
	Update(ctx context.Context, currentVersion string) (io.ReadCloser, error)
}

// Update updates the current binary with the chosen resolver
func Update(resolver Resolver, currentVersion string) (bool, error) {
	return UpdateWithContext(context.TODO(), resolver, currentVersion)
}

// UpdateWithContext updates the current binary with the chosen resolver and accepts a context
func UpdateWithContext(ctx context.Context, resolver Resolver, currentVersion string) (bool, error) {
	reader, err := resolver.Update(ctx, currentVersion)
	if err != nil {
		return false, err
	}
	defer reader.Close()

	githubUpdater := updater.New()
	err = githubUpdater.Apply(reader)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			return true, nil
		}
		return false, errors.Wrap(err, "could not apply the update")
	}

	return true, nil
}
