package user

import (
	"effective_mobile/internal/domain"
)

type Service interface {
	Index(page int, name, surname string) ([]*domain.User, error)
	Store(name, surname string, patronymic *string) (*domain.User, error)
	Show(id string) (*domain.User, error)
	Update(id, name, surname string, patronymic, gender, countryId *string, age *int) (*domain.User, error)
	Delete(id string) (string, error)
}
