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
	Name      string    `json:"name"`
	Count     int       `json:"count"`
	Countries []Country `json:"country"`
}

type Country struct {
	Name        *string `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Api struct {
	domain string
}

func New(domain string) *Api {
	strings.TrimRight(domain, "/")

	return &Api{domain: domain}
}

func (a *Api) NatByName(name string) (*string, error) {
	const op = "repository.api.nationalize.user.NationalityByName"

	var nat *string

	params := url.Values{}
	params.Add("name", name)

	fullUrl := fmt.Sprintf("%s/?%s", a.domain, params.Encode())

	response, err := http.Get(fullUrl)
	if err != nil {
		return nat, fmt.Errorf("%s: %w", op, err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nat, fmt.Errorf("%s: %w", op, err)
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nat, fmt.Errorf("%s: %w", op, err)
	}

	var maxP float64
	var country Country

	for _, item := range user.Countries {
		if item.Probability > maxP {
			maxP = item.Probability
			country = item
		}
	}

	return country.Name, nil
}
