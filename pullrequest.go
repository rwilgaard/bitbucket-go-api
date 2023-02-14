package gobitbucket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type ToRef struct {
    ID           string      `json:"id"`
    Type         string      `json:"type"`
    LatestCommit string      `json:"latestCommit"`
    DisplayID    string      `json:"displayId"`
    Repository   *Repository `json:"repository"`
}

type FromRef struct {
    ID           string      `json:"id"`
    Type         string      `json:"type"`
    LatestCommit string      `json:"latestCommit"`
    DisplayID    string      `json:"displayId"`
    Repository   *Repository `json:"repository"`
}

type Participant struct {
    Role               string `json:"role"`
    User               *User  `json:"user"`
    Status             string `json:"status"`
    LastReviewedCommit string `json:"lastReviewdCommit"`
    Approved           bool   `json:"approved"`
}

type PullRequest struct {
    Version      int32          `json:"version"`
    Locked       bool           `json:"locked"`
    ID           int64          `json:"id"`
    State        string         `json:"state"`
    Open         bool           `json:"open"`
    Title        string         `json:"title"`
    Closed       bool           `json:"closed"`
    ToRef        *ToRef         `json:"toRef"`
    CreatedDate  int64          `json:"createdDate"`
    FromRef      *FromRef       `json:"fromRef"`
    Participants []*Participant `json:"participants"`
    ClosedDate   int64          `json:"closedDate"`
    Reviewers    []*Participant `json:"reviewers"`
    Description  string         `json:"description"`
    UpdatedDate  int64          `json:"updatedDate"`
}

type PullRequestList struct {
    Values        []*PullRequest `json:"values"`
    Size          int            `json:"size"`
    Limit         int            `json:"limit"`
    Start         int32          `json:"start"`
    IsLastPage    bool           `json:"isLastPage"`
    NextPageStart int32          `json:"nextPageStart"`
}

type PullRequestsQuery struct {
    ProjectKey     string
    RepositorySlug string
    WithAttributes string // (optional) defaults to true, whether to return additional pull request attributes
    At             string // (optional) a fully-qualified branch ID to find pull requests to or from, such as refs/heads/master
    WithProperties string // (optional) defaults to true, whether to return additional pull request properties
    FilterText     string // (optional) If specified, only pull requests where the title or description contains the supplied string will be returned.
    State          string // (optional, defaults to OPEN). Supply ALL to return pull request in any state. If a state is supplied only pull requests in the specified state will be returned. Either OPEN, DECLINED or MERGED.
    Order          string // (optional, defaults to NEWEST) the order to return pull requests in, either OLDEST (as in: "oldest first") or NEWEST.
    Direction      string // (optional, defaults to INCOMING) the direction relative to the specified repository. Either INCOMING or OUTGOING.
    Start          int    // Start number for the page (inclusive). If not passed, first page is assumed.
    Limit          int    // Number of items to return. If not passed, a page size of 25 is used.
}

func addPullRequestsQueryParams(query PullRequestsQuery) *url.Values {
    data := url.Values{}
    if query.WithAttributes != "" {
        data.Set("withAttributes", query.WithAttributes)
    }
    if query.At != "" {
        data.Set("at", query.At)
    }
    if query.WithProperties != "" {
        data.Set("withProperties", query.WithProperties)
    }
    if query.FilterText != "" {
        data.Set("filterText", query.FilterText)
    }
    if query.State != "" {
        data.Set("state", query.State)
    }
    if query.Order != "" {
        data.Set("order", query.Order)
    }
    if query.Direction != "" {
        data.Set("direction", query.Direction)
    }
    if query.Start != 0 {
        data.Set("start", strconv.Itoa(query.Start))
    }
    if query.Limit != 0 {
        data.Set("limit", strconv.Itoa(query.Limit))
    }
    return &data
}

func (a *API) GetPullRequests(query PullRequestsQuery) (*PullRequestList, error) {
    p := fmt.Sprintf("/rest/api/latest/projects/%s/repos/%s/pull-requests", query.ProjectKey, query.RepositorySlug)
    u, err := url.ParseRequestURI(a.endpoint.String() + p)
    if err != nil {
        return nil, err
    }
    u.RawQuery = addPullRequestsQueryParams(query).Encode()
    req, err := http.NewRequest("GET", u.String(), nil)
    if err != nil {
        return nil, err
    }
    req.SetBasicAuth(a.username, a.token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := a.Client.Do(req)
    if err != nil {
        return nil, err
    }

    res, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if err = resp.Body.Close(); err != nil {
        return nil, err
    }

    var pr PullRequestList
    if err := json.Unmarshal(res, &pr); err != nil {
        return nil, err
    }

    return &pr, nil
}
