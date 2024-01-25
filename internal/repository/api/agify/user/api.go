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
	Age   *int   `json:"age"`
	Count int    `json:"count"`
	Name  string `json:"name"`
}

type Api struct {
	domain string
}

func New(domain string) *Api {
	strings.TrimRight(domain, "/")

	return &Api{domain: domain}
}

func (a *Api) AgeByName(name string) (*int, error) {
	const op = "repository.api.agify.user.AgeByName"

	var age *int = nil

	params := url.Values{}
	params.Add("name", name)

	fullUrl := fmt.Sprintf("%s/?%s", a.domain, params.Encode())

	response, err := http.Get(fullUrl)
	if err != nil {
		return age, fmt.Errorf("%s: %w", op, err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return age, fmt.Errorf("%s: %w", op, err)
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return age, fmt.Errorf("%s: %w", op, err)
	}

	return user.Age, nil
}
