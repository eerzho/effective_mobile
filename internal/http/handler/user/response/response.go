package response

import (
	"effective_mobile/internal/domain"
	"effective_mobile/internal/lib/api/response"
)

type Index struct {
	response.Success
	Users []*domain.User `json:"users"`
}

type Store struct {
	response.Success
	User *domain.User `json:"user"`
}

type Show struct {
	response.Success
	User *domain.User `json:"user"`
}

type Update struct {
	response.Success
	User *domain.User `json:"user"`
}

type Delete struct {
	response.Success
	Id string `json:"id"`
}
