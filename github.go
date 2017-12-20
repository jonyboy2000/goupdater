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
	Token        string
	Repo         string
	Owner        string
}

// NewGithub creates a new instance of Github
func NewGithub(token string, owner string, repo string) (*Github, error) {
	if token == "" {
		return nil, ErrMissingGithubToken
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	githubClient := github.NewClient(tc)

	return &Github{
		GithubClient: githubClient,
		HTTPClient:   http.DefaultClient,
		Token:        token,
		Owner:        owner,
		Repo:         repo,
	}, nil
}

// Update checks and returns a new binary on github releases
func (p *Github) Update(currentVersion string) (io.ReadCloser, error) {
	ctx := context.Background()

	r, _, err := p.GithubClient.Repositories.GetLatestRelease(ctx, p.Owner, p.Repo)
	if err != nil {
		return nopCloser{}, errors.Wrap(err, "could not fetch github release")
	}

	if *r.TagName != currentVersion {
		downloadURL, err := p.getPlatformReleaseURL(r)
		if err != nil {
			return nopCloser{}, err
		}

		q := downloadURL.Query()
		q.Add("access_token", p.Token)
		downloadURL.RawQuery = q.Encode()

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
