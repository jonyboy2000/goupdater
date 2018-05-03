package goupdater

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var (
	// ErrMissingGithubToken is used when a github token is not provided
	ErrMissingGithubToken = errors.New("to check for updates you must provide a github token")
)

// Github represents a github releases resolver
type Github struct {
	githubClient *github.Client
	httpClient   *http.Client

	token string
	owner string
	repo  string
}

// NewGithub creates a new instance of Github
func NewGithub(opts ...Option) (*Github, error) {
	return NewGithubWithContext(context.TODO(), opts...)
}

// NewGithubWithContext creates a new instance of Github and accepts a context param
func NewGithubWithContext(ctx context.Context, opts ...Option) (*Github, error) {
	var ts oauth2.TokenSource
	updater := Github{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(&updater)
	}

	if updater.token != "" {
		ts = oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: updater.token},
		)
	}

	tc := oauth2.NewClient(ctx, ts)
	updater.githubClient = github.NewClient(tc)

	return &updater, nil
}

// Update checks and returns a new binary on github releases
func (p *Github) Update(ctx context.Context, currentVersion string) (io.ReadCloser, error) {
	nopReader := bytes.NewReader(nil)
	nopCloser := ioutil.NopCloser(nopReader)

	r, _, err := p.githubClient.Repositories.GetLatestRelease(ctx, p.owner, p.repo)
	if err != nil {
		return nopCloser, errors.Wrap(err, "could not fetch github release")
	}

	if *r.TagName != currentVersion {
		downloadURL, err := p.getPlatformReleaseURL(r)
		if err != nil {
			return nopCloser, err
		}

		if p.token != "" {
			q := downloadURL.Query()
			q.Add("access_token", p.token)
			downloadURL.RawQuery = q.Encode()
		}

		req, err := http.NewRequest(http.MethodGet, downloadURL.String(), nil)
		if err != nil {
			return nopCloser, errors.Wrap(err, "could not create a request for the release download URL")
		}

		req.Header.Add("Accept", "application/octet-stream")

		resp, err := p.httpClient.Do(req)
		if err != nil {
			return nopCloser, errors.Wrap(err, "could not make the request for a release")
		}

		return resp.Body, nil
	}

	return nopCloser, nil
}

func (p *Github) getPlatformReleaseURL(r *github.RepositoryRelease) (*url.URL, error) {
	for _, asset := range r.Assets {
		if strings.Contains(asset.GetName(), runtime.GOOS) {
			return url.Parse(asset.GetURL())
		}
	}

	return nil, errors.New("could not find a valid URL")
}
