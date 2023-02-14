package gobitbucket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Tag struct {
    Hash            string `json:"hash"`
    LatestCommit    string `json:"latestCommit"`
    LatestChangeset string `json:"latestChangeset"`
    ID              string `json:"id"`
    Type            string `json:"type"`
    DisplayID       string `json:"displayId"`
}

type TagList struct {
    Values        []*Tag `json:"values"`
    Size          int    `json:"size"`
    Limit         int    `json:"limit"`
    Start         int32  `json:"start"`
    IsLastPage    bool   `json:"isLastPage"`
    NextPageStart int32  `json:"nextPageStart"`
}

type TagsQuery struct {
    ProjectKey     string
    RepositorySlug string
    OrderBy        string // Ordering of refs either ALPHABETICAL (by name) or MODIFICATION (last updated)
    FilterText     string
    Start          int
    Limit          int
}

func addTagsQueryParams(query TagsQuery) *url.Values {
    data := url.Values{}
    if query.OrderBy != "" {
        data.Set("orderBy", query.OrderBy)
    }
    if query.FilterText != "" {
        data.Set("filterText", query.FilterText)
    }
    if query.Start != 0 {
        data.Set("start", strconv.Itoa(query.Start))
    }
    if query.Limit != 0 {
        data.Set("limit", strconv.Itoa(query.Limit))
    }
    return &data
}

func (a *API) GetTags(query TagsQuery) (*TagList, error) {
    p := fmt.Sprintf("/rest/api/latest/projects/%s/repos/%s/tags", query.ProjectKey, query.RepositorySlug)
    u, err := url.ParseRequestURI(a.endpoint.String() + p)
    if err != nil {
        return nil, err
    }
    u.RawQuery = addTagsQueryParams(query).Encode()
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

    var tags TagList
    if err := json.Unmarshal(res, &tags); err != nil {
        return nil, err
    }

    return &tags, nil
}
