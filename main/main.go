package main

import (
	"fmt"

	gobitbucket "github.com/rwilgaard/bitbucket-go-api"
)


func main() {
    api, err := gobitbucket.NewAPI("", "", "")
    if err != nil {
        panic(err)
    }

    query := gobitbucket.RepositoriesQuery{
        Limit: 9999,
    }
    repos, err := api.GetRepositories(query)
    if err != nil {
        panic(err)
    }

    for _, r := range repos.Values {
        fmt.Printf("%+v\n", r)
    }
}

