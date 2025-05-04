package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stepan41k/Effective-Mobile/internal/domain/models"
)

const (
	host               = "localhost:8082"
	pageSize           = 10
	page               = 1
	invalidName        = "Abcdefghijklmnopqrstwxzqwertyasdfgh"
	invalidSurname     = "Abcdefghijklmnopqrstwxzqwertyasdfgh"
	invalidPatronymic  = "Abcdefghijklmnopqrstwxzqwertyasdfgh"
	invalidAge         = 200
	invalidGender      = "Panzerkampfwagen"
	invalidNationalize = "Some Nationalize"
)

func TestMobileCreate_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	name := gofakeit.FirstName()
	surname := gofakeit.LastName()

	e.POST("/profile/new").
		WithJSON(models.NewPerson{
			Name:    name,
			Surname: surname,
		}).
		Expect().
		Status(http.StatusOK)
}

func TestMobileUpdate_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	name := gofakeit.FirstName()
	surname := gofakeit.LastName()

	guid := e.POST("/profile/new").
		WithJSON(models.NewPerson{
			Name:    name,
			Surname: surname,
		}).Expect().
		JSON().
		Object().
		Value("data").
		String().Raw()

	e.PATCH("/profile/update").
		WithJSON(models.UpdatedPerson{
			GUID:        guid,
			Name:        gofakeit.FirstName(),
			Surname:     gofakeit.LastName(),
			Patronymic:  gofakeit.MiddleName(),
			Age:         gofakeit.Number(10, 80),
			Gender:      gofakeit.Gender(),
			Nationalize: gofakeit.Country(),
		}).
		Expect().
		Status(http.StatusOK)
}

func TestMobileGet_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	name := gofakeit.FirstName()
	surname := gofakeit.LastName()
	patronymic := gofakeit.MiddleName()
	gender := gofakeit.Gender()
	nationalize := gofakeit.Country()

	guid := e.POST("/profile/new").
		WithJSON(models.NewPerson{
			Name:    name,
			Surname: surname,
		}).Expect().
		JSON().
		Object().
		Value("data").
		String().Raw()

	e.PATCH("/profile/update").
		WithJSON(models.UpdatedPerson{
			GUID:        guid,
			Patronymic:  patronymic,
			Gender:      gender,
			Nationalize: nationalize,
		}).
		Expect().
		Status(http.StatusOK)

	e.GET("/profile/profiles").
		WithJSON(models.GetPerson{
			Name:        name,
			Surname:     surname,
			Patronymic:  patronymic,
			Gender:      gender,
			Nationalize: nationalize,
			PageSize:    pageSize,
			Page:        page,
		}).
		Expect().
		Status(http.StatusOK)
}

func TestMobileDelete_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	name := gofakeit.FirstName()
	surname := gofakeit.LastName()

	guid := e.POST("/profile/new").
		WithJSON(models.NewPerson{
			Name:    name,
			Surname: surname,
		}).Expect().
		JSON().
		Object().
		Value("data").
		String().Raw()

	e.DELETE("/profile/delete").
		WithJSON(models.DeletePerson{
			GUID: guid,
		}).
		Expect().
		Status(http.StatusOK)
}

func TestCreate_FailCases(t *testing.T) {
	cases := []struct {
		title     string
		name      string
		surname   string
		respError string
	}{
		{
			title:     "Create profile with empty name",
			name:      "",
			surname:   gofakeit.LastName(),
			respError: "field Name is a required field",
		},
		{
			title:     "Create profile with empty surname",
			name:      gofakeit.FirstName(),
			surname:   "",
			respError: "field Surname is a required field",
		},
		{
			title:     "Create profile with too large name",
			name:      invalidName,
			surname:   gofakeit.LastName(),
			respError: "field Name must have less than 20 characters",
		},
		{
			title:     "Create profile with too large surname",
			name:      gofakeit.FirstName(),
			surname:   invalidSurname,
			respError: "field Surname must have less than 30 characters",
		},
	}

	for _, tt := range cases {

		t.Run(tt.title, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/profile/new").
				WithJSON(models.NewPerson{
					Name:    tt.name,
					Surname: tt.surname,
				}).Expect().JSON().Object()

			if tt.respError != "" {
				resp.NotContainsKey("data")

				resp.Value("error").String().IsEqual(tt.respError)

				return
			}
		})
	}
}

func TestUpdate_FailCases(t *testing.T) {
	cases := []struct {
		title       string
		name        string
		surname     string
		patronymic  string
		age         int
		gender      string
		nationalize string
		respError   string
	}{
		{
			title:       "Update without GUID",
			name:        gofakeit.FirstName(),
			surname:     gofakeit.LastName(),
			patronymic:  gofakeit.MiddleName(),
			age:         gofakeit.Number(10, 80),
			gender:      "RU",
			nationalize: gofakeit.Country(),
			respError:   "field GUID is a required field",
		},
		{
			title:       "Update with invalid name",
			name:        invalidName,
			surname:     gofakeit.LastName(),
			patronymic:  gofakeit.MiddleName(),
			age:         gofakeit.Number(10, 80),
			gender:      gofakeit.Gender(),
			nationalize: "RU",
			respError:   "field Name must have less than 20 characters",
		},
		{
			title:       "Update with invalid surname",
			name:        gofakeit.FirstName(),
			surname:     invalidSurname,
			patronymic:  gofakeit.MiddleName(),
			age:         gofakeit.Number(10, 80),
			gender:      gofakeit.Gender(),
			nationalize: "RU",
			respError:   "field Surname must have less than 30 characters",
		},
		{
			title:       "Update with invalid patronymic",
			name:        gofakeit.FirstName(),
			surname:     gofakeit.LastName(),
			patronymic:  invalidPatronymic,
			age:         gofakeit.Number(10, 80),
			gender:      gofakeit.Gender(),
			nationalize: "RU",
			respError:   "field Patronymic must have less than 25 characters",
		},
		{
			title:       "Update with invalid age",
			name:        gofakeit.FirstName(),
			surname:     gofakeit.LastName(),
			patronymic:  gofakeit.MiddleName(),
			age:         invalidAge,
			gender:      gofakeit.Gender(),
			nationalize: "RU",
			respError:   "field Age must have less than 130",
		},
		{
			title:       "Update with invalid gender",
			name:        gofakeit.FirstName(),
			surname:     gofakeit.LastName(),
			patronymic:  gofakeit.MiddleName(),
			age:         gofakeit.Number(10, 80),
			gender:      invalidGender,
			nationalize: "RU",
			respError:   "field Gender must have less than 6 characters",
		},
		{
			title:       "Update with invalid nationalize",
			name:        gofakeit.FirstName(),
			surname:     gofakeit.LastName(),
			patronymic:  gofakeit.MiddleName(),
			age:         gofakeit.Number(10, 80),
			gender:      gofakeit.Gender(),
			nationalize: invalidNationalize,
			respError:   "field Nationalize must have less than 3 characters",
		},
	}

	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			guid := e.POST("/profile/new").
				WithJSON(models.NewPerson{
					Name:    gofakeit.FirstName(),
					Surname: gofakeit.LastName(),
				}).Expect().
				JSON().
				Object().
				Value("data").
				String().Raw()

			if tt.title == "Update without GUID" {
				resp := e.PATCH("/profile/update").
					WithJSON(models.UpdatedPerson{
						GUID: "",
					}).Expect().JSON().Object()

				if tt.respError != "" {
					resp.NotContainsKey("data")

					resp.Value("error").String().IsEqual(tt.respError)

					return
				}
			}

			resp := e.PATCH("/profile/update").
				WithJSON(models.UpdatedPerson{
					GUID: guid,
					Name: tt.name,
					Surname: tt.surname,
					Patronymic: tt.patronymic,
					Age: tt.age,
					Gender: tt.gender,
					Nationalize: tt.nationalize,
				}).Expect().JSON().Object()

			if tt.respError != "" {
				resp.NotContainsKey("data")

				resp.Value("error").String().IsEqual(tt.respError)

				return
			}
		})
	}
}

func TestGet_FailCases(t *testing.T) {
	cases := []struct {
		title    string
		page     int
		pageSize int
		respError string
	}{
		{
			title:       "Get profiles without number of page",
			pageSize:     5,
			respError: "field Page is a required field",
		},
		{
			title:       "Get profiles without size of page",
			page:     5,
			respError: "field PageSize is a required field",
		},
	}

	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.GET("/profile/profiles").
				WithJSON(models.GetPerson{
					Page:    tt.page,
					PageSize: tt.pageSize,
				}).Expect().JSON().Object()

			if tt.respError != "" {
				resp.NotContainsKey("data")

				resp.Value("error").String().IsEqual(tt.respError)

				return
			}
		})
	}
}

func TestDelete_FailCases(t *testing.T) {
	cases := []struct {
		title     string
		guid      string
		respError string
	}{
		{
			title:     "Delete profile with empty GUID",
			guid:      "",
			respError: "field GUID is a required field",
		},
		{
			title:     "Delete non-existent profile",
			guid:      "random-non-existed-guid",
			respError: "profile not found",
		},
	}

	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.DELETE("/profile/delete").
				WithJSON(models.DeletePerson{
					GUID: tt.guid,
				}).Expect().JSON().Object()

			if tt.respError != "" {
				resp.NotContainsKey("data")

				resp.Value("error").String().IsEqual(tt.respError)

				return
			}
		})
	}
}
