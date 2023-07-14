package gobitbucket

import (
    "net/http"
    "net/url"
    "strconv"
)

type Project struct {
    Key         string `json:"key"`
    Name        string `json:"name,omitempty"`
    ID          uint   `json:"id,omitempty"`
    Type        string `json:"type,omitempty"`
    Public      bool   `json:"public,omitempty"`
    Scope       string `json:"scope,omitempty"`
    Description string `json:"description,omitempty"`
    Namespace   string `json:"namespace,omitempty"`
    Avatar      string `json:"avatar,omitempty"`
}

type ProjectList struct {
    Values []*Project `json:"values"`
    Page
}

type ProjectsQuery struct {
    Name       string // Name to filter by.
    Permission string // Permission to filter by
    Start      uint   // Start number for the page (inclusive). If not passed, first page is assumed.
    Limit      uint   // Number of items to return. If not passed, a page size of 25 is used.
}

func getProjectsQueryParams(query ProjectsQuery) *url.Values {
    data := url.Values{}
    if query.Name != "" {
        data.Set("name", query.Name)
    }
    if query.Permission != "" {
        data.Set("permission", query.Permission)
    }
    if query.Start != 0 {
        data.Set("start", strconv.FormatUint(uint64(query.Start), 10))
    }
    if query.Limit != 0 {
        data.Set("limit", strconv.FormatUint(uint64(query.Limit), 10))
    }
    return &data
}

func (a *API) GetProjects(query ProjectsQuery) (*ProjectList, *http.Response, error) {
    p := getProjectsQueryParams(query)
    req, err := a.NewRequest("GET", "/rest/api/latest/projects", nil, p)
    if err != nil {
        return nil, nil, err
    }

    projects := ProjectList{
        Page: Page{
            IsLastPage: true,
        },
    }
    resp, err := a.Do(req, &projects)
    if err != nil {
        return nil, resp, err
    }

    return &projects, resp, nil
}
