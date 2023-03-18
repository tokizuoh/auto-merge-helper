package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	repository := os.Getenv("GITHUB_REPOSITORY")
	sha := os.Getenv("GITHUB_SHA")
	if token == "" || repository == "" || sha == "" {
		log.Fatalln("failed to retrieve expected environment variable")
	}

	ownerRepositoryName := strings.Split(repository, "/")
	if len(ownerRepositoryName) != 2 {
		log.Fatalf("failed to split the \"%s\" as expected.", repository)
	}
	owner := ownerRepositoryName[0]
	repositoryName := ownerRepositoryName[1]

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

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

	var query struct {
		Repository struct {
			Object struct {
				Commit struct {
					StatusCheckRollup struct {
						Contexts struct {
							Nodes []Context
						} `graphql:"contexts(first: 100)"`
					}
					AbbreviatedOid string
				} `graphql:"... on Commit"`
			} `graphql:"object(expression: $expression)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"name":       githubv4.String(repositoryName),
		"expression": githubv4.String(sha),
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Fatal(err)
	}

	f := true
	for _, node := range query.Repository.Object.Commit.StatusCheckRollup.Contexts.Nodes {
		if node.CheckRun.Name != "" {
			if node.CheckRun.Name == "auto-merge-helper" {
				continue
			}
			if node.CheckRun.Conclusion != githubv4.CheckConclusionStateSuccess {
				f = false
			}
		} else if node.StatusContext.Context != "" {
			if node.StatusContext.State != githubv4.StatusStateSuccess {
				f = false
			}
		} else {
			log.Fatal("invalid inline fragment: does not match the expected type")
		}
	}

	if f {
		log.Println("all success")
	} else {
		log.Fatalln("failed to run CI successfully, at least one CI has failed")
	}
}
