package gobitbucket

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
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
    Values        []*Commit `json:"values"`
    Size          int       `json:"size"`
    Limit         int       `json:"limit"`
    Start         int32     `json:"start"`
    IsLastPage    bool      `json:"isLastPage"`
    NextPageStart int32     `json:"nextPageStart"`
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
    Start          int
    Limit          int
}

func (a *API) getCommitsEndpoint(query CommitsQuery) (*url.URL, error) {
    p := fmt.Sprintf("/rest/api/latest/projects/%s/repos/%s/commits", query.ProjectKey, query.RepositorySlug)
    return url.ParseRequestURI(a.endpoint.String() + p)
}

func addCommitsQueryParams(query CommitsQuery) *url.Values {
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
        data.Set("start", strconv.Itoa(query.Start))
    }
    if query.Limit != 0 {
        data.Set("limit", strconv.Itoa(query.Limit))
    }
    return &data
}

func (a *API) GetCommits(query CommitsQuery) (*CommitList, error) {
    u, err := a.getCommitsEndpoint(query)
    if err != nil {
        return nil, err
    }
    u.RawQuery = addCommitsQueryParams(query).Encode()
    req, err := http.NewRequest("GET", u.String(), nil)
    req.SetBasicAuth(a.username, a.token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := a.Client.Do(req)
    if err != nil {
        panic(err)
    }

    res, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    err = resp.Body.Close()
    if err != nil {
        panic(err)
    }

    var commits CommitList
    json.Unmarshal(res, &commits)

    return &commits, nil
}
