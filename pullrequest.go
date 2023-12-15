package gobitbucket

import (
    "fmt"
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
    Version      int32          `json:"version,omitempty"`
    Locked       bool           `json:"locked,omitempty"`
    ID           int64          `json:"id,omitempty"`
    State        string         `json:"state,omitempty"`
    Open         bool           `json:"open,omitempty"`
    Title        string         `json:"title"`
    Closed       bool           `json:"closed,omitempty"`
    ToRef        *ToRef         `json:"toRef,omitempty"`
    CreatedDate  int64          `json:"createdDate,omitempty"`
    FromRef      *FromRef       `json:"fromRef,omitempty"`
    Participants []*Participant `json:"participants,omitempty"`
    ClosedDate   int64          `json:"closedDate,omitempty"`
    Reviewers    []*Participant `json:"reviewers,omitempty"`
    Description  string         `json:"description,omitempty"`
    UpdatedDate  int64          `json:"updatedDate,omitempty"`
}

type PullRequestList struct {
    Values []*PullRequest `json:"values"`
    Page
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
    Start          uint   // Start number for the page (inclusive). If not passed, first page is assumed.
    Limit          uint   // Number of items to return. If not passed, a page size of 25 is used.
}

type InboxPullRequestCount struct {
    Count uint `json:"count"`
}

func getPullRequestsQueryParams(query PullRequestsQuery) *url.Values {
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
        data.Set("start", strconv.FormatUint(uint64(query.Start), 10))
    }
    if query.Limit != 0 {
        data.Set("limit", strconv.FormatUint(uint64(query.Limit), 10))
    }
    return &data
}

func (a *API) GetPullRequests(query PullRequestsQuery) (*PullRequestList, *http.Response, error) {
    params := getPullRequestsQueryParams(query)
    path := fmt.Sprintf("/rest/api/latest/projects/%s/repos/%s/pull-requests", query.ProjectKey, query.RepositorySlug)
    req, err := a.NewRequest("GET", path, nil, params)
    if err != nil {
        return nil, nil, err
    }

    pr := PullRequestList{
        Page: Page{
            IsLastPage: true,
        },
    }
    resp, err := a.Do(req, &pr)
    if err != nil {
        return nil, resp, err
    }

    return &pr, resp, nil
}

func (a *API) GetInboxPullRequestCount() (*InboxPullRequestCount, *http.Response, error) {
    path := "/rest/api/latest/inbox/pull-requests/count"
    req, err := a.NewRequest("GET", path, nil, nil)
    if err != nil {
        return nil, nil, err
    }

    pr := new(InboxPullRequestCount)
    resp, err := a.Do(req, pr)
    if err != nil {
        return nil, resp, err
    }

    return pr, resp, nil
}

