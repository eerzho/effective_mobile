package domain

type User struct {
	Id         string  `json:"id"`
	Age        *int    `json:"age"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic"`
	Gender     *string `json:"gender"`
	CountryId  *string `json:"country_id"`
}
