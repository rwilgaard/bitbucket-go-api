package gobitbucket

import (
	"fmt"
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

func getTagsQueryParams(query TagsQuery) *url.Values {
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

func (a *API) GetTags(query TagsQuery) (*TagList, *http.Response, error) {
    params := getTagsQueryParams(query)
    path := fmt.Sprintf("/rest/api/latest/projects/%s/repos/%s/tags", query.ProjectKey, query.RepositorySlug)
    req, err := a.NewRequest("GET", path, nil, params)
    if err != nil {
        return nil, nil, err
    }

    tags := TagList{
        IsLastPage: true,
    }
    resp, err := a.Do(req, &tags) 
    if err != nil {
        return nil, resp, err
    }

    return &tags, resp, nil
}
