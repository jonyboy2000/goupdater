package goupdater_test

import (
	"context"
	"fmt"

	"github.com/italolelis/goupdater"
)

func Example_privateRepo() {
	// we define all the necessary configurations.
	var (
		currentVersion = "v1.2.3"    // the version of your app
		githubToken    = "yourToken" // This is only required for private github repositories
		owner          = "hellofresh"
		repo           = "myRepo"
	)

	// We choose the github release resolver. You can implement your own resolver by implementing
	// the goupdater.Resolver interface
	resolver, err := goupdater.NewGithub(goupdater.GithubOpts{
		Token: githubToken,
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		panic(err)
	}

	// Updates the running binary to the latest available version
	updated, err := goupdater.Update(resolver, currentVersion)
	if err != nil {
		panic(err)
	}

	if updated {
		fmt.Print("You are now using the latest version!")
	} else {
		fmt.Print("You already have the latest version!")
	}
}

func Example_publicRepoWithContext() {
	// we define all the necessary configurations.
	var (
		currentVersion = "v1.2.3" // the version of your app
		owner          = "hellofresh"
		repo           = "myRepo"
	)

	// we create a context to be used, this could also be a ContextTimeout that would
	// fail in case of a network error
	ctx := context.Background()

	// We choose the github release resolver. You can implement your own resolver by implementing
	// the goupdater.Resolver interface
	resolver, err := goupdater.NewGithubWithContext(ctx, goupdater.GithubOpts{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		panic(err)
	}

	// Updates the running binary to the latest available version
	updated, err := goupdater.UpdateWithContext(ctx, resolver, currentVersion)
	if err != nil {
		panic(err)
	}

	if updated {
		fmt.Print("You are now using the latest version!")
	} else {
		fmt.Print("You already have the latest version!")
	}
}
