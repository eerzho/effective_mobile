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
	const op = "repository.nationalize.user.ByName"

	var maxP float64
	var maxC *string

	params := url.Values{}
	params.Add("name", name)

	fullUrl := fmt.Sprintf("%s/?%s", r.domain, params.Encode())

	response, err := http.Get(fullUrl)
	if err != nil {
		return maxC, fmt.Errorf("%s: %w", op, err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return maxC, fmt.Errorf("%s: %w", op, err)
	}

	var user struct {
		Name      string `json:"name"`
		Count     int    `json:"count"`
		Countries []struct {
			CountryId   *string `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return maxC, fmt.Errorf("%s: %w", op, err)
	}

	for _, item := range user.Countries {
		if item.Probability >= maxP {
			maxP = item.Probability
			maxC = item.CountryId
		}
	}

	return maxC, nil
}
