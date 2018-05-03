package goupdater

// Option represents the client options
type Option func(*Github)

// WithToken sets the github token
func WithToken(token string) Option {
	return func(g *Github) {
		g.token = token
	}
}

// WithOwner sets the github owner
func WithOwner(owner string) Option {
	return func(g *Github) {
		g.owner = owner
	}
}

// WithRepo sets the github repo
func WithRepo(repo string) Option {
	return func(g *Github) {
		g.repo = repo
	}
}
