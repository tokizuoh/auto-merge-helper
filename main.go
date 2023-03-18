package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	body := []byte(``)
	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", buf)
	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("no such environment variable: GITHUB_TOKEN")
	}

	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// TODO
	log.Println(res)
}
