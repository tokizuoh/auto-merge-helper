package main

import (
	"log"
	"os"
	"strings"
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

	query, err := fetch(token, owner, repositoryName, sha)
	if err != nil {
		log.Fatalln(err)
	}

	ok, err := checkAllSuccess(query.Repository.Object.Commit.StatusCheckRollup.Contexts.Nodes)
	if err != nil {
		log.Fatalln(err)
	}

	if ok {
		log.Println("all success")
	} else {
		log.Fatalln("failed to run CI successfully, at least one CI has failed")
	}
}
