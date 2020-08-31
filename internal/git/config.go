package git

import (
	"fmt"
	"strings"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
)

// ExtractRepoFromConfig gets the repo name from the Git config.
func ExtractRepoFromConfig(ctx *context.Context) (result config.Repo, err error) {
	if !IsRepo(ctx) {
		return result, ErrNotRepository{ctx.CurrentDirectory}
	}
	out, err := Run(ctx, "config", "--get", "remote.origin.url")
	if err != nil {
		return result, fmt.Errorf("repository doesn't have an `origin` remote")
	}
	return ExtractRepoFromURL(out), nil
}

// ExtractRepoFromURL gets the repo name from the remote URL.
func ExtractRepoFromURL(s string) config.Repo {
	// removes the .git suffix and any new lines
	s = strings.NewReplacer(
		".git", "",
		"\n", "",
	).Replace(s)
	// if the URL contains a :, indicating a SSH config,
	// remove all chars until it, including itself
	// on HTTP and HTTPS URLs it will remove the http(s): prefix,
	// which is ok. On SSH URLs the whole user@server will be removed,
	// which is required.
	s = s[strings.LastIndex(s, ":")+1:]
	// split by /, the last to parts should be the owner and name
	ss := strings.Split(s, "/")
	return config.Repo{
		Owner: ss[len(ss)-2],
		Name:  ss[len(ss)-1],
	}
}
