package user

import (
	"fmt"
	"log/slog"
	"sync"

	"effective_mobile/internal/domain"
	"effective_mobile/internal/lib/logger/sl"
)

type Service struct {
	log           *slog.Logger
	repo          Repo
	agerRepo      AgerRepo
	genderRepo    GenderRepo
	countryerRepo CountryerRepo
}

func New(
	log *slog.Logger,
	repo Repo,
	agerRepo AgerRepo,
	genderRepo GenderRepo,
	countryerRepo CountryerRepo,
) *Service {
	return &Service{
		log:           log,
		repo:          repo,
		agerRepo:      agerRepo,
		genderRepo:    genderRepo,
		countryerRepo: countryerRepo,
	}
}

func (s *Service) Index(page, size int, name, surname, patronymic, gender, countryId string, age int) ([]*domain.User, error) {
	const op = "service.user.Index"

	log := s.log.With(slog.String("op", op))

	log.Info("getting users")

	list, err := s.repo.List(page, size, name, surname, patronymic, gender, countryId, age)
	if err != nil {
		log.Error("failed to get users")

		return []*domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got users")

	return list, nil
}

func (s *Service) Store(name, surname string, patronymic *string) (*domain.User, error) {
	const op = "service.user.Store"

	log := s.log.With(slog.String("op", op))

	log.Info("saving user")

	var wg sync.WaitGroup
	wg.Add(3)

	var age *int
	go func() {
		defer wg.Done()
		log.Info("getting user age")
		var err error
		age, err = s.agerRepo.ByName(name)
		if err != nil {
			log.Error("failed to getting age", sl.Err(err))
		}
	}()

	var gender *string
	go func() {
		defer wg.Done()
		log.Info("getting user gender")
		var err error
		gender, err = s.genderRepo.ByName(name)
		if err != nil {
			log.Error("failed to getting gender", sl.Err(err))
		}
	}()

	var countryId *string
	go func() {
		defer wg.Done()
		log.Info("getting user country_id")
		var err error
		countryId, err = s.countryerRepo.ByName(name)
		if err != nil {
			log.Error("failed to getting country_id", sl.Err(err))
		}
	}()

	wg.Wait()

	user, err := s.repo.Save(name, surname, patronymic, gender, countryId, age)
	if err != nil {
		log.Error("failed to save user")

		return &domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("saved user")

	return user, nil
}

func (s *Service) Show(id string) (*domain.User, error) {
	const op = "service.user.Show"

	log := s.log.With(slog.String("op", op))

	log.Info("getting user")

	user, err := s.repo.GetById(id)
	if err != nil {
		log.Error("failed to get user by id")

		return &domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got user")

	return user, nil
}

func (s *Service) Update(id, name, surname string, patronymic, gender, countryId *string, age *int) (*domain.User, error) {
	const op = "service.user.Update"

	log := s.log.With(slog.String("op", op))

	log.Info("updating user")

	user, err := s.repo.Update(id, name, surname, patronymic, gender, countryId, age)
	if err != nil {
		log.Error("failed to update user by id")

		return &domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user updated")

	return user, nil
}

func (s *Service) Delete(id string) (string, error) {
	const op = "service.user.delete"

	log := s.log.With(slog.String("op", op))

	log.Info("deleting user")

	id, err := s.repo.DelById(id)
	if err != nil {
		log.Error("failed to delete user by id")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user deleted")

	return id, nil
}
