package gobitbucket

type User struct {
    Name         string `json:"name"`
    ID           int32  `json:"id"`
    Type         string `json:"type"`
    DisplayName  string `json:"displayName"`
    Slug         string `json:"slug"`
    Active       bool   `json:"active"`
    EmailAddress string `json:"emailAddress"`
}
