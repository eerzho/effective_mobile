package user

import "effective_mobile/internal/domain"

type Repo interface {
	List(page, size int, name, surname, patronymic, gender, countryId string, age int) ([]*domain.User, error)
	GetById(id string) (*domain.User, error)
	Exists(id string) (bool, error)
	Save(name, surname string, patronymic, gender, countryId *string, age *int) (*domain.User, error)
	Update(id, name, surname string, patronymic, gender, countryId *string, age *int) (*domain.User, error)
	DelById(id string) (string, error)
}

type AgerRepo interface {
	ByName(name string) (*int, error)
}

type GenderRepo interface {
	ByName(name string) (*string, error)
}

type CountryerRepo interface {
	ByName(name string) (*string, error)
}
