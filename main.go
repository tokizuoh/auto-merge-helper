package main

import (
	"context"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	repository := os.Getenv("GITHUB_REPOSITORY")
	sha := os.Getenv("GITHUB_SHA")
	if token == "" || repository == "" || sha == "" {
		log.Fatalln("not exist expected environment variable")
	}

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
			} `graphql:"object(expression: \"8cab9ceabc5e0af9a6d407b80357a\")"`
		} `graphql:"repository(owner: \"tokizuoh\", name: \"citrus\")"`
	}
	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		log.Fatal(err)
	}

	f := true
	for _, node := range query.Repository.Object.Commit.StatusCheckRollup.Contexts.Nodes {
		if node.CheckRun.Name != "" {
			if node.CheckRun.Conclusion != githubv4.CheckConclusionStateSuccess {
				f = false
			}
		} else if node.StatusContext.Context != "" {
			if node.StatusContext.State != githubv4.StatusStateSuccess {
				f = false
			}
		} else {
			log.Fatal("failed to expand inline fragment")
		}
	}

	if f {
		log.Println("all success")
	} else {
		log.Fatalln("TODO")
	}
}
