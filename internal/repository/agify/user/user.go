package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Repo struct {
	domain string
}

func New(domain string) *Repo {
	strings.TrimRight(domain, "/")

	return &Repo{domain: domain}
}

func (r *Repo) ByName(name string) (*int, error) {
	const op = "repository.agify.user.ByName"

	var user struct {
		Age   *int   `json:"age"`
		Count int    `json:"count"`
		Name  string `json:"name"`
	}

	params := url.Values{}
	params.Add("name", name)

	fullUrl := fmt.Sprintf("%s/?%s", r.domain, params.Encode())

	response, err := http.Get(fullUrl)
	if err != nil {
		return user.Age, fmt.Errorf("%s: %w", op, err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return user.Age, fmt.Errorf("%s: %w", op, err)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user.Age, fmt.Errorf("%s: %w", op, err)
	}

	return user.Age, nil
}
