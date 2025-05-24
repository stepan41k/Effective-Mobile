package models

type Person struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationalize string `json:"nationalize,omitempty"`
}

type NewPerson struct {
	GUID        string `json:"guid,omitempty" validate:"omitempty"`
	Name        string `json:"name" validate:"required,min=1,max=20" example:"Igor"`
	Surname     string `json:"surname" validate:"required,min=1,max=30" example:"Zaycev"`
	Patronymic  string `json:"patronymic,omitempty" validate:"omitempty,min=1,max=25" example:"Vladimirovich"`
}

type EnrichedPerson struct {
	GUID        string `json:"guid"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationalize string `json:"nationalize"`
}

type GetPerson struct {
	Name        string `json:"name,omitempty" example:"John"`
	Surname     string `json:"surname,omitempty" example:"Wick"`
	Patronymic  string `json:"patronymic,omitempty" example:"Ivanovich"`
	Age         int    `json:"age,omitempty" example:"28"`
	Greater     bool   `json:"greater,omitempty" example:"true"`
	Gender      string `json:"gender,omitempty" example:"male"`
	Nationalize string `json:"nationalize,omitempty" example:"US"`
	PageSize    int    `json:"page_size" validate:"required" example:"10"`
	Page        int    `json:"page" validate:"required" example:"3"`
}

type UpdatedPerson struct {
	GUID        string `json:"guid" validate:"required" example:"3EWQbnsu-2!IHY389-ewqh312"`
	Name        string `json:"new_name,omitempty" validate:"omitempty,min=1,max=20" example:"Valeriy"`
	Surname     string `json:"new_surname,omitempty" validate:"omitempty,min=1,max=30" example:"Popov"`
	Patronymic  string `json:"patronymic,omitempty" validate:"omitempty,min=1,max=25" example:"Valentinovich"`
	Age         int    `json:"age,omitempty" validate:"omitempty,gte=0,lte=130" example:"33"`
	Gender      string `json:"gender,omitempty" validate:"omitempty,max=6" example:"male"`
	Nationalize string `json:"nationalize,omitempty" validate:"omitempty,max=3" example:"RU"`
}

type DeletePerson struct {
	GUID string `json:"guid" validate:"required" example:"ewqehQWE231u-Snu3h21sj-321s"`
}
