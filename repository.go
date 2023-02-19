package gobitbucket

import (
	"net/http"
	"net/url"
	"strconv"
)

type Repository struct {
    Name        string            `json:"name"`
    ID          int               `json:"id,omitempty"`
    Slug        string            `json:"slug,omitempty"`
    Public      bool              `json:"public,omitempty"`
    Archived    bool              `json:"archived,omitempty"`
    Description string            `json:"description,omitempty"`
    State       string            `json:"state,omitempty"`
    Project     Project           `json:"project,omitempty"`
    Links       map[string][]Link `json:"links,omitempty"`
}

type Link struct {
    Name string `json:"name,omitempty"`
    Href string `json:"href"`
}

type RepositoryList struct {
    Values        []*Repository `json:"values"`
    Size          int           `json:"size"`
    Limit         int           `json:"limit"`
    Start         int32         `json:"start"`
    IsLastPage    bool          `json:"isLastPage"`
    NextPageStart int32         `json:"nextPageStart"`
}

type RepositoriesQuery struct {
    Archived    string // ACTIVE,ARCHIVED OR ALL. Default is ACTIVE
    ProjectName string
    ProjectKey  string
    Visibility  string // public,private
    Name        string
    Permission  string // REPO_READ,REPO_WRITE,REPO_ADMIN
    State       string // AVAILABLE,INITIALISING,INITIALISATION_FAILED
    Start       int
    Limit       int
}

func getRepositoriesQueryParams(query RepositoriesQuery) *url.Values {
    data := url.Values{}
    if query.Archived != "" {
        data.Set("archived", query.Archived)
    }
    if query.ProjectName != "" {
        data.Set("projectname", query.ProjectName)
    }
    if query.ProjectKey != "" {
        data.Set("projectkey", query.ProjectKey)
    }
    if query.Visibility != "" {
        data.Set("visibility", query.Visibility)
    }
    if query.Name != "" {
        data.Set("name", query.Name)
    }
    if query.Permission != "" {
        data.Set("permission", query.Permission)
    }
    if query.State != "" {
        data.Set("state", query.State)
    }
    if query.Start != 0 {
        data.Set("start", strconv.Itoa(query.Start))
    }
    if query.Limit != 0 {
        data.Set("limit", strconv.Itoa(query.Limit))
    }
    return &data
}

func (a *API) GetRepositories(query RepositoriesQuery) (*RepositoryList, *http.Response, error) {
    p := getRepositoriesQueryParams(query)
    req, err := a.NewRequest("GET", "/rest/api/latest/repos", nil, p)
    if err != nil {
        return nil, nil, err
    }

    repos := RepositoryList{
        IsLastPage: true,
    }
    resp, err := a.Do(req, &repos)
    if err != nil {
        return nil, resp, err
    }

    return &repos, resp, nil
}
