package main

import (
	"context"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
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

	// TODO: Handle
	for _, n := range query.Repository.Object.Commit.StatusCheckRollup.Contexts.Nodes {
		// log.Println(n.CheckRun.Name)
		// log.Println(n.CheckRun.Conclusion)
		log.Println(n.CheckRun.DetailsUrl)
		// log.Println(n.StatusContext.Context)
		// log.Println(n.StatusContext.State)
		log.Println(n.StatusContext.TargetUrl)
		log.Println("---")
	}
}
