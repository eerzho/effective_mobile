package request

type Store struct {
	Name       string  `json:"name" validate:"required,min=2,max=255"`
	Surname    string  `json:"surname" validate:"required,min=2,max=255"`
	Patronymic *string `json:"patronymic,omitempty" validate:"omitempty,min=2,max=255"`
}

type Update struct {
	Age        *int    `json:"age,omitempty" validate:"omitempty"`
	Name       string  `json:"name" validate:"required,min=2,max=255"`
	Surname    string  `json:"surname" validate:"required,min=2,max=255"`
	Patronymic *string `json:"patronymic,omitempty" validate:"omitempty,min=2,max=255"`
	Gender     *string `json:"gender,omitempty" validate:"omitempty,min=2,max=255"`
	CountryId  *string `json:"country_id,omitempty" validate:"omitempty,min=2,max=255"`
}
