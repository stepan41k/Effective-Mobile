package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/stepan41k/Effective-Mobile/internal/domain/models"
	"github.com/stepan41k/Effective-Mobile/internal/lib/api/logger/sl"
	resp "github.com/stepan41k/Effective-Mobile/internal/lib/api/response"
	"github.com/stepan41k/Effective-Mobile/internal/service"
)

type Profile interface {
	TakeProfiles(ctx context.Context, profile models.GetPerson) (profiles []models.Person, err error)
	RemoveProfile(ctx context.Context, profile models.DeletePerson) (guid []byte, err error)
	UpdateProfile(ctx context.Context, profile models.UpdatedPerson) (guid []byte, err error)
	NewProfile(ctx context.Context, profile models.EnrichedPerson) (guid []byte, err error)
}

type ProfileHandler struct {
	profile Profile
	log     *slog.Logger
}

func New(profile Profile, log *slog.Logger) *ProfileHandler {
	return &ProfileHandler{
		profile: profile,
		log:     log,
	}
}

const (
	EnvAge = "AGE_API"
	EnvGender = "GENDER_API"
	EnvNationalize = "NATIONALIZE_API"
)

// @Summary Get
// @Tags profile
// @Description Accepts filters and outputs profiles based on them
// @ID get-profiles
// @Accept  json
// @Produce  json
// @Param input body models.GetPerson true "page and size of page is necessary"
// @Success 200 {object} response.SuccessResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /get [post]
func (m *ProfileHandler) TakeProfiles(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.profile.TakeProfile"

		log := m.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.GetPerson

		err := render.Decode(r, &req)
		flag := CheckForErrors(req, w, r, log, err)
		if flag {
			return
		}

		profiles, err := m.profile.TakeProfiles(ctx, req)
		if err != nil {
			if errors.Is(err, service.ErrProfilesNotFound) {
				log.Warn("profiles not found")

				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusConflict,
					Error:  "profiles not found",
				})
			}

			log.Error("internal error")

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "internal error",
			})

			return
		}

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data:   profiles,
		})
	}
}

// @Summary Delete
// @Tags profile
// @Description Accepts profile GUID and remove this profile
// @ID delete-profile
// @Accept  json
// @Produce  json
// @Param input body models.DeletePerson true "GUID is necessary"
// @Success 200 {object} response.SuccessResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /delete [delete]
func (m *ProfileHandler) RemoveProfile(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.music.RemoveProfile"

		log := m.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.DeletePerson

		err := render.Decode(r, &req)
		flag := CheckForErrors(req, w, r, log, err)
		if flag {
			return
		}

		guid, err := m.profile.RemoveProfile(ctx, req)
		if err != nil {
			if errors.Is(err, service.ErrProfileNotFound) {
				log.Warn("profile not found")

				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusConflict,
					Error:  "profile not found",
				})

				return
			}
			log.Error("internal error")

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "internal error",
			})

			return
		}

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data:   string(guid),
		})
	}
}

// @Summary Update
// @Tags profile
// @Description Accepts profile GUID and remove this profile
// @ID update-profile
// @Accept  json
// @Produce  json
// @Param input body models.UpdatedPerson true "GUID is necessary"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /update [patch]
func (m *ProfileHandler) UpdateProfile(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.music.UpdateProfile"

		log := m.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.UpdatedPerson

		err := render.Decode(r, &req)
		flag := CheckForErrors(req, w, r, log, err)
		if flag {
			return
		}

		guid, err := m.profile.UpdateProfile(ctx, req)
		if err != nil {
			if errors.Is(err, service.ErrProfileNotFound) {
				log.Warn("profile not found")

				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusConflict,
					Error:  "profile not found",
				})

				return
			}

			if errors.Is(err, service.ErrNoChanges) {
				log.Warn("nothing to update or profile not exists")

				render.Status(r, http.StatusConflict)

				render.JSON(w, r, resp.ErrorResponse{
					Status: http.StatusConflict,
					Error:  "nothing to update or profile not exists",
				})

				return
			}

			log.Error("internal error", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "internal error",
			})

			return
		}

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data:   string(guid),
		})
	}
}

// @Summary Create
// @Tags profile
// @Description Accepts name, surname and patronymic and creates profile
// @ID create-profile
// @Accept  json
// @Produce  json
// @Param input body models.NewPerson true "name and surname is necessary"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /create [post]
func (m *ProfileHandler) NewProfile(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.music.NewProfile"

		log := m.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.NewPerson
		var profile models.EnrichedPerson

		err := render.Decode(r, &req)
		flag := CheckForErrors(req, w, r, log, err)
		if flag {
			return
		}

		var age models.Age
		ageReq, err := http.Get(os.Getenv(EnvAge) + req.Name)
		if err != nil {
			log.Error("failed to get age")

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "failed to get age",
			})

			return
		}

		render.DecodeJSON(ageReq.Body, &age)
		profile.Age = age.Age

		var gender models.Gender
		genderReq, err := http.Get(os.Getenv(EnvGender) + req.Name)
		if err != nil {
			log.Error("failed to get gender")

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "failed to get gender",
			})

			return
		}

		render.DecodeJSON(genderReq.Body, &gender)
		profile.Gender = gender.Gender

		var nationalize models.Nationalize
		nationalizeReq, err := http.Get(os.Getenv(EnvNationalize) + req.Name)
		if err != nil {
			log.Error("failed to get nationalize")

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "failed to get nationalize",
			})

			return
		}

		render.DecodeJSON(nationalizeReq.Body, &nationalize)
		profile.Nationalize = nationalize.Country[0].CountryID

		profile.GUID = req.GUID
		profile.Name = req.Name
		profile.Surname = req.Surname
		profile.Patronymic = req.Patronymic

		log.Info(profile.Name, profile.Surname, profile.Patronymic, profile.Gender, profile.Nationalize)

		guid, err := m.profile.NewProfile(ctx, profile)
		if err != nil {
			log.Error("internal error", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "internal error",
			})

			return
		}

		render.JSON(w, r, resp.SuccessResponse{
			Status: http.StatusOK,
			Data:   string(guid),
		})
	}
}

func CheckForErrors(req any, w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) bool {

	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.Status(r, http.StatusConflict)

			render.JSON(w, r, resp.ErrorResponse{
				Status: http.StatusConflict,
				Error:  "empty request",
			})

			return true
		}

		log.Error("failed to decode request")

		render.Status(r, http.StatusBadRequest)

		render.JSON(w, r, resp.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "failed to decode request",
		})
		return true
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)

		log.Error("invalid request", sl.Err(err))

		render.JSON(w, r, resp.ValidationError(validateErr))

		return true
	}

	return false
}
