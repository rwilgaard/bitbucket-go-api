package gobitbucket

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
)

type API struct {
    endpoint *url.URL
    Client   *http.Client
    username string
    token    string
}

func NewAPI(endpoint string, username string, token string) (*API, error) {
    if len(endpoint) == 0 {
        return nil, errors.New("url empty")
    }

    u, err := url.ParseRequestURI(endpoint)
    if err != nil {
        return nil, err
    }

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
    }

    a := new(API)
    a.endpoint = u
    a.Client = &http.Client{Transport: tr}
    a.username = username
    a.token = token

    return a, nil
}
