package domain

type User struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic"`
	Sex         *string `json:"sex"`
	Nationality *string `json:"nationality"`
	Age         *int    `json:"age"`
}
