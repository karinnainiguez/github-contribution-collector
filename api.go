package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"reflect"

	"github.com/google/go-querystring/query"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func newClient() *github.Client {
	ctx := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUBKEY")},
	)

	tokenClient := oauth2.NewClient(ctx, tokenService)

	client := github.NewClient(tokenClient)
	return client
}

func getRepo(c *github.Client, owner string, name string) *github.Repository {
	ctx := context.Background()
	repo, _, err := c.Repositories.Get(ctx, owner, name)

	handle(err)
	return repo
}

// manual api call, sdk does not provide working functionality
func getRepos(c *github.Client, org string) []*github.Repository {

	ctx := context.Background()

	opts := &github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allRepos []*github.Repository
	for {
		u := fmt.Sprintf("orgs/%v/repos", org)
		u, err := addOptions(u, opts)
		handle(err)

		req, err := c.NewRequest("GET", u, nil)
		handle(err)

		var repos []*github.Repository
		resp, err := c.Do(ctx, req, &repos)
		handle(err)

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allRepos
}

// helper method to add options to manual api call
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func getIssues(c *github.Client, repo *github.Repository, user string) []*github.Issue {
	opts := &github.IssueListByRepoOptions{
		State:       "all",
		Creator:     user,
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allIssues []*github.Issue
	for {
		issues, resp, err := c.Issues.ListByRepo(
			context.Background(),
			repo.GetOwner().GetLogin(),
			repo.GetName(),
			opts,
		)
		handle(err)

		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allIssues
}
