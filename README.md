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
var currentVersion = "v1.2.3"
var githubToken = "yourToken"

// We choose the github release provider
resolver, err := goupdater.NewGithub(githubToken, "hellofresh", "myRepo")
failOnError(err, "could not create github resolver")

// Updates the running binary to the latest available version
updated, err := goupdater.Update(resolver, currentVersion)
failOnError(err, "could not update binary")

if updated {
    log.Info("You are now using the latest version!")
} else {
    log.Info("You already have the latest version!")
}
```
