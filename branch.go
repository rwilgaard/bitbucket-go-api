package gobitbucket

import (
    "fmt"
    "net/http"
    "net/url"
    "strconv"
)

type Branch struct {
    Default         bool   `json:"default,omitempty"`
    Type            any    `json:"type,omitempty"`
    DisplayID       string `json:"displayId,omitempty"`
    LatestCommit    string `json:"latestCommit,omitempty"`
    LatestChangeset string `json:"latestChangeset,omitempty"`
    ID              string `json:"id"`
}

type BranchList struct {
    Values []*Branch `json:"values"`
    Page
}

type BranchesQuery struct {
    ProjectKey     string
    RepositorySlug string
    BoostMatches   bool   // Controls whether exact and prefix matches will be boosted to the top
    OrderBy        string // Ordering of refs either ALPHABETICAL (by name) or MODIFICATION (last updated)
    Details        bool   // Whether to retrieve plugin-provided metadata about each branch
    FilterText     string // The text to match on
    Base           string // Base branch or tag to compare each branch to (for the metadata providers that uses that information
    Start          uint
    Limit          uint
}

func getBranchesQueryParams(query BranchesQuery) *url.Values {
    data := url.Values{}
    if query.BoostMatches {
        data.Set("boostMatches", "true")
    }
    if query.OrderBy != "" {
        data.Set("orderBy", query.OrderBy)
    }
    if query.Details {
        data.Set("details", "true")
    }
    if query.FilterText != "" {
        data.Set("filterText", query.FilterText)
    }
    if query.Base != "" {
        data.Set("base", query.Base)
    }
    if query.Start != 0 {
        data.Set("start", strconv.FormatUint(uint64(query.Start), 10))
    }
    if query.Limit != 0 {
        data.Set("limit", strconv.FormatUint(uint64(query.Limit), 10))
    }
    return &data
}

func (a *API) GetBranches(query BranchesQuery) (*BranchList, *http.Response, error) {
    params := getBranchesQueryParams(query)
    path := fmt.Sprintf("/rest/api/latest/projects/%s/repos/%s/branches", query.ProjectKey, query.RepositorySlug)
    req, err := a.NewRequest("GET", path, nil, params)
    if err != nil {
        return nil, nil, err
    }

    branches := BranchList{
        Page: Page{
            IsLastPage: true,
        },
    }
    resp, err := a.Do(req, &branches)
    if err != nil {
        return nil, resp, err
    }

    return &branches, resp, err
}
