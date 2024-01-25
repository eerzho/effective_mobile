package request

type Store struct {
	Name       string  `json:"name" validate:"required,min=2,max=255"`
	Surname    string  `json:"surname" validate:"required,min=2,max=255"`
	Patronymic *string `json:"patronymic,omitempty" validate:"omitempty,min=2,max=255"`
}

type Update struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Surname     string  `json:"surname" validate:"required,min=2,max=255"`
	Patronymic  *string `json:"patronymic,omitempty" validate:"omitempty,min=2,max=255"`
	Sex         *string `json:"sex,omitempty" validate:"omitempty,min=2,max=255"`
	Nationality *string `json:"nationality,omitempty" validate:"omitempty,min=2,max=255"`
	Age         *int    `json:"age,omitempty" validate:"omitempty"`
}
