package gobitbucket

import (
    "fmt"
    "net/http"
    "net/url"
    "strconv"
)

type Comitter struct {
    Name         string `json:"name,omitempty"`
    EmailAddress string `json:"emailAddress"`
}

type Author struct {
    Name         string `json:"name,omitempty"`
    EmailAddress string `json:"emailAddress"`
}

type Commit struct {
    Message            string    `json:"message"`
    Author             *Author   `json:"author"`
    AuthorTimestamp    int64     `json:"authorTimestamp"`
    Committer          *Comitter `json:"committer"`
    CommitterTimestamp int64     `json:"committerTimestamp"`
    ID                 string    `json:"id"`
    DisplayID          string    `json:"displayId"`
}

type CommitList struct {
    Values []*Commit `json:"values"`
    Page
}

type CommitsQuery struct {
    ProjectKey     string
    RepositorySlug string
    Path           string // An optional path to filter commits by
    WithCounts     int    // Optionally include the total number of commits and total number of unique authors
    FollowRenames  string // If true, the commit history of the specified file will be followed past renames. Only valid for a path to a single file.
    Until          string // The commit ID (SHA1) or ref (inclusively) to retrieve commits before
    Since          string // The commit ID or ref (exclusively) to retrieve commits after
    Merges         string // exclude,include,only for merge commits
    IgnoreMissing  string
    Start          uint
    Limit          uint
}

func getCommitsQueryParams(query CommitsQuery) *url.Values {
    data := url.Values{}
    if query.Path != "" {
        data.Set("path", query.Path)
    }
    if query.WithCounts != 0 {
        data.Set("withCounts", strconv.Itoa(query.WithCounts))
    }
    if query.FollowRenames != "" {
        data.Set("followRenames", query.FollowRenames)
    }
    if query.Until != "" {
        data.Set("until", query.Until)
    }
    if query.Since != "" {
        data.Set("since", query.Since)
    }
    if query.Merges != "" {
        data.Set("merges", query.Merges)
    }
    if query.IgnoreMissing != "" {
        data.Set("ignoreMissing", query.IgnoreMissing)
    }
    if query.Start != 0 {
        data.Set("start", strconv.FormatUint(uint64(query.Start), 10))
    }
    if query.Limit != 0 {
        data.Set("limit", strconv.FormatUint(uint64(query.Limit), 10))
    }
    return &data
}

func (a *API) GetCommits(query CommitsQuery) (*CommitList, *http.Response, error) {
    params := getCommitsQueryParams(query)
    path := fmt.Sprintf("/rest/api/latest/projects/%s/repos/%s/commits", query.ProjectKey, query.RepositorySlug)
    req, err := a.NewRequest("GET", path, nil, params)
    if err != nil {
        return nil, nil, err
    }

    commits := CommitList{
        Page: Page{
            IsLastPage: true,
        },
    }
    resp, err := a.Do(req, &commits)
    if err != nil {
        return nil, resp, err
    }

    return &commits, resp, nil
}
