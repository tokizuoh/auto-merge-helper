package main

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Context struct {
	CheckRun struct {
		Conclusion githubv4.CheckConclusionState
		DetailsUrl githubv4.URI
		Name       string
	} `graphql:"... on CheckRun"`
	StatusContext struct {
		State     githubv4.StatusState
		TargetUrl githubv4.URI
		Context   string
	} `graphql:"... on StatusContext"`
}

type Query struct {
	Repository struct {
		Object struct {
			Commit struct {
				StatusCheckRollup struct {
					Contexts struct {
						Nodes []Context
					} `graphql:"contexts(first: 100)"`
				}
			} `graphql:"... on Commit"`
		} `graphql:"object(expression: $expression)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func checkAllSuccess(ctxs []Context) (bool, error) {
	f := true
	for _, ctx := range ctxs {
		if ctx.CheckRun.Name != "" {
			if ctx.CheckRun.Name == "auto-merge-helper" {
				continue
			}
			if ctx.CheckRun.Conclusion != githubv4.CheckConclusionStateSuccess {
				f = false
			}
		} else if ctx.StatusContext.Context != "" {
			if ctx.StatusContext.State != githubv4.StatusStateSuccess {
				f = false
			}
		} else {
			return false, fmt.Errorf("invalid inline fragment: does not match the expected type")
		}
	}

	return f, nil
}

func fetch(token, owner, repoName, sha string) (*Query, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	query := Query{}
	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"name":       githubv4.String(repoName),
		"expression": githubv4.String(sha),
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}

	return &query, nil
}
