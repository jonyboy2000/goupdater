package goupdater

import (
	"context"
	"io"
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
	GithubClient *github.Client
	HTTPClient   *http.Client
	Opts         GithubOpts
}

// GithubOpts represents the available github resolver options
type GithubOpts struct {
	Token string
	Owner string
	Repo  string
}

// NewGithub creates a new instance of Github
func NewGithub(opts GithubOpts) (*Github, error) {
	return NewGithubWithContext(context.TODO(), opts)
}

// NewGithubWithContext creates a new instance of Github and accepts a context param
func NewGithubWithContext(ctx context.Context, opts GithubOpts) (*Github, error) {
	var ts oauth2.TokenSource

	if opts.Token != "" {
		ts = oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: opts.Token},
		)
	}

	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	return &Github{
		GithubClient: githubClient,
		HTTPClient:   http.DefaultClient,
		Opts:         opts,
	}, nil
}

// Update checks and returns a new binary on github releases
func (p *Github) Update(ctx context.Context, currentVersion string) (io.ReadCloser, error) {
	r, _, err := p.GithubClient.Repositories.GetLatestRelease(ctx, p.Opts.Owner, p.Opts.Repo)
	if err != nil {
		return nopCloser{}, errors.Wrap(err, "could not fetch github release")
	}

	if *r.TagName != currentVersion {
		downloadURL, err := p.getPlatformReleaseURL(r)
		if err != nil {
			return nopCloser{}, err
		}

		if p.Opts.Token != "" {
			q := downloadURL.Query()
			q.Add("access_token", p.Opts.Token)
			downloadURL.RawQuery = q.Encode()
		}

		req, err := http.NewRequest(http.MethodGet, downloadURL.String(), nil)
		if err != nil {
			return nopCloser{}, errors.Wrap(err, "could not create a request for the release download URL")
		}

		req.Header.Add("Accept", "application/octet-stream")

		resp, err := p.HTTPClient.Do(req)
		if err != nil {
			return nopCloser{}, errors.Wrap(err, "could not make the request for a release")
		}

		return resp.Body, nil
	}

	return nopCloser{}, nil
}

func (p *Github) getPlatformReleaseURL(r *github.RepositoryRelease) (*url.URL, error) {
	for _, asset := range r.Assets {
		if strings.Contains(asset.GetName(), runtime.GOOS) {
			return url.Parse(asset.GetURL())
		}
	}

	return nil, errors.New("could not find a valid URL")
}
