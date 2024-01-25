package user

import "effective_mobile/internal/domain"

type Repository interface {
	List(page int, name, surname string) ([]domain.User, error)
	GetById(id string) (domain.User, error)
	Exists(id string) (bool, error)
	Save(name, surname string, patronymic, sex, nationality *string, age *int) (domain.User, error)
	Update(id, name, surname string, patronymic, sex, nationality *string, age *int) (domain.User, error)
	DelById(id string) (string, error)
}
