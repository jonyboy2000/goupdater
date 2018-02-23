# Goupdater

> Easily build self-updating programs

Goupdater makes it easier for you to update your Go programs (or other single-file targets). A program can update itself by replacing its executable file with a new version.

It provides the flexibility to implement different updating user experiences like auto-updating, or manual user-initiated updates.

## Install

We use [dep](https://github.com/golang/dep) to manage our dependencies. you can easily add this library to your application by doing:

Using `dep`

```sh
dep ensure --add github.com/italolelis/goupdater
```

Using `go get`

```sh
go get github.com/italolelis/goupdater
```

## Usage

Goupdater allows you to have different resolvers from where to fetch your binaries. At the moment we only support `github releases`, but you can easily write your own resolver.

Here is an example of how you can update your application using the `github releases` resover:

```go
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
```

## Contributing
To start contributing, please check [CONTRIBUTING](CONTRIBUTING).

## Documentation
Go Docs: [godoc.org/github.com/italolelis/goupdater](https://godoc.org/github.com/italolelis/goupdater)
