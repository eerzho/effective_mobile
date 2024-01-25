package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type User struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Sex         *string `json:"gender"`
	Probability float64 `json:"probability"`
}

type Api struct {
	domain string
}

func New(domain string) *Api {
	strings.TrimRight(domain, "/")

	return &Api{domain: domain}
}

func (a *Api) SexByName(name string) (*string, error) {
	const op = "repository.api.genderize.user.SexByName"

	var sex *string = nil

	params := url.Values{}
	params.Add("name", name)

	fullUrl := fmt.Sprintf("%s/?%s", a.domain, params.Encode())

	response, err := http.Get(fullUrl)
	if err != nil {
		return sex, fmt.Errorf("%s: %w", op, err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return sex, fmt.Errorf("%s: %w", op, err)
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return sex, fmt.Errorf("%s: %w", op, err)
	}

	return user.Sex, nil
}
