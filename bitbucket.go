package gobitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ErrorResponse struct {
	*http.Response
	Errors []ErrorMessage
}

type ErrorMessage struct {
	Message   string `json:"message,omitempty"`
	Exception string `json:"exceptionName,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		e.Response.Request.Method, e.Response.Request.URL,
		e.Response.StatusCode, e.Errors)
}

func CheckResponse(res *http.Response) error {
	if c := res.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	return &ErrorResponse{Response: res}
}

func (a *API) NewRequest(method string, path string, body interface{}, params *url.Values) (*http.Request, error) {
    u, err := url.ParseRequestURI(a.endpoint.String() + path)
    if err != nil {
        return nil, err
    }

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

    u.RawQuery = params.Encode()

    req, err := http.NewRequest(method, u.String(), buf)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.SetBasicAuth(a.username, a.token)

    return req, nil
}

func (a *API) Do(req *http.Request, i interface{}) (*http.Response, error) {
    resp, err := a.Client.Do(req)
    if err != nil {
        return nil, err
    }

    if err = CheckResponse(resp); err != nil {
        return resp, err
    }

    if i == nil {
        return resp, nil
    }

    if err = json.NewDecoder(resp.Body).Decode(i); err != nil {
        return resp, err
    }

    if err = resp.Body.Close(); err != nil {
        return resp, err
    }

    return resp, nil
}
