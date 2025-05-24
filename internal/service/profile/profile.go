package music

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/stepan41k/Effective-Mobile/internal/domain/models"
	"github.com/stepan41k/Effective-Mobile/internal/lib/api/logger/sl"
	"github.com/stepan41k/Effective-Mobile/internal/service"
	"github.com/stepan41k/Effective-Mobile/internal/storage"
)

type Profile interface {
	TakeProfiles(ctx context.Context, person models.GetPerson) (persons []models.Person, err error)
	RemoveProfile(ctx context.Context, person models.DeletePerson) (guid []byte, err error)
	UpdateProfile(ctx context.Context, person models.UpdatedPerson) (guid []byte, err error)
	NewProfile(ctx context.Context, person models.EnrichedPerson) (guid []byte, err error)
}

type ProfileService struct {
	profile Profile
	log     *slog.Logger
}

func New(profile Profile, log *slog.Logger) *ProfileService {
	return &ProfileService{
		profile: profile,
		log:     log,
	}
}

func (m *ProfileService) TakeProfiles(ctx context.Context, person models.GetPerson) ([]models.Person, error) {
	const op = "service.music.GetProfiles"

	log := m.log.With(
		slog.String("op", op),
		slog.String("name", person.Name),
		slog.String("surname", person.Surname),
	)

	log.Info("getting profiles")

	profiles, err := m.profile.TakeProfiles(ctx, person)
	if err != nil {
		if errors.Is(err, storage.ErrProfilesNotFound) {
			log.Warn("profiles not found")

			return nil, fmt.Errorf("%s: %w", op, service.ErrProfilesNotFound)
		}

		log.Error("failed to get profiles", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got profiles")

	return profiles, nil
}

func (m *ProfileService) RemoveProfile(ctx context.Context, person models.DeletePerson) ([]byte, error) {
	const op = "service.music.DeleteProfile"

	log := m.log.With(
		slog.String("op", op),
		slog.String("guid", person.GUID),
	)

	log.Info("deleting profile")

	guid, err := m.profile.RemoveProfile(ctx, person)
	if err != nil {
		if errors.Is(err, storage.ErrProfileNotFound) {
			log.Warn("profile not found")

			return nil, fmt.Errorf("%s: %w", op, service.ErrProfileNotFound)
		}
		log.Error("failed to delete profile", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("profile deleted")

	return guid, nil
}

func (m *ProfileService) UpdateProfile(ctx context.Context, person models.UpdatedPerson) ([]byte, error) {
	const op = "service.profile.UpdateProfile"

	log := m.log.With(
		slog.String("op", op),
		slog.String("guid", person.GUID),
	)

	log.Info("updating profile")

	id, err := m.profile.UpdateProfile(ctx, person)
	if err != nil {
		if errors.Is(err, storage.ErrProfileNotFound) {
			log.Warn("profile not found")

			return nil, fmt.Errorf("%s: %w", op, service.ErrProfileNotFound)
		}

		if errors.Is(err, storage.ErrNoChanges) {
			log.Warn("nothing to update")

			return nil, fmt.Errorf("%s: %w", op, service.ErrNoChanges)
		}
		log.Error("failed to update profile", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("profile updated")

	return id, nil
}

func (m *ProfileService) NewProfile(ctx context.Context, person models.EnrichedPerson) ([]byte, error) {
	const op = "service.music.NewProfile"

	log := m.log.With(
		slog.String("op", op),
		slog.String("name", person.Name),
		slog.String("surname", person.Surname),
	)

	log.Info("creating new profile")

	guid, err := uuid.NewRandom()
	if err != nil {
		log.Error("failed to generate guid")

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	person.GUID = guid.String()

	id, err := m.profile.NewProfile(ctx, person)
	if err != nil {
		log.Error("failed to add profile", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("profile added")

	return id, nil
}
