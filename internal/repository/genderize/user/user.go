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

func (r *Repo) ByName(name string) (*string, error) {
	const op = "repository.genderize.user.ByName"

	var user struct {
		Count       int     `json:"count"`
		Name        string  `json:"name"`
		Gender      *string `json:"gender"`
		Probability float64 `json:"probability"`
	}

	params := url.Values{}
	params.Add("name", name)

	fullUrl := fmt.Sprintf("%s/?%s", r.domain, params.Encode())

	response, err := http.Get(fullUrl)
	if err != nil {
		return user.Gender, fmt.Errorf("%s: %w", op, err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return user.Gender, fmt.Errorf("%s: %w", op, err)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user.Gender, fmt.Errorf("%s: %w", op, err)
	}

	return user.Gender, nil
}
