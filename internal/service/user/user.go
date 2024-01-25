package user

import (
	"fmt"
	"log/slog"
	"sync"

	"effective_mobile/internal/domain"
	"effective_mobile/internal/lib/logger/sl"
)

type User struct {
	log           *slog.Logger
	repository    Repository
	ageRepository AgeRepository
	sexRepository SexRepository
	natRepository NatRepository
}

func New(
	log *slog.Logger,
	repository Repository,
	ageRepository AgeRepository,
	sexRepository SexRepository,
	natRepository NatRepository,
) *User {
	return &User{
		log:           log,
		repository:    repository,
		ageRepository: ageRepository,
		sexRepository: sexRepository,
		natRepository: natRepository,
	}
}

func (u *User) Index(page int, name, surname string) ([]domain.User, error) {
	const op = "service.user.Index"

	log := u.log.With(slog.String("op", op))

	log.Info("getting users")

	list, err := u.repository.List(page, name, surname)
	if err != nil {
		log.Error("failed to get users")

		return []domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got users")

	return list, nil
}

func (u *User) Store(name, surname string, patronymic *string) (domain.User, error) {
	const op = "service.user.Store"

	log := u.log.With(slog.String("op", op))

	log.Info("saving user")

	var temp struct {
		age *int
		sex *string
		nat *string
	}

	var wg sync.WaitGroup
	mutex := sync.Mutex{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		log.Info("getting user age")
		age, err := u.ageRepository.AgeByName(name)
		if err != nil {
			log.Error("failed to getting age", sl.Err(err))
		}

		mutex.Lock()
		temp.age = age
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()

		log.Info("getting user sex")
		sex, err := u.sexRepository.SexByName(name)
		if err != nil {
			log.Error("failed to getting sex", sl.Err(err))
		}

		mutex.Lock()
		temp.sex = sex
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()

		log.Info("getting user nationality")
		nat, err := u.natRepository.NatByName(name)
		if err != nil {
			log.Error("failed to getting nationality", sl.Err(err))
		}

		mutex.Lock()
		temp.nat = nat
		mutex.Unlock()
	}()

	wg.Wait()

	user, err := u.repository.Save(name, surname, patronymic, temp.sex, temp.nat, temp.age)
	if err != nil {
		log.Error("failed to save user")

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("saved user")

	return user, nil
}

func (u *User) Show(id string) (domain.User, error) {
	const op = "service.user.Show"

	log := u.log.With(slog.String("op", op))

	log.Info("getting user")

	user, err := u.repository.GetById(id)
	if err != nil {
		log.Error("failed to get user by id")

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got user")

	return user, nil
}

func (u *User) Update(id, name, surname string, patronymic, sex, nationality *string, age *int) (domain.User, error) {
	const op = "service.user.Update"

	log := u.log.With(slog.String("op", op))

	log.Info("updating user")

	user, err := u.repository.Update(id, name, surname, patronymic, sex, nationality, age)
	if err != nil {
		log.Error("failed to update user by id")

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user updated")

	return user, nil
}

func (u *User) Delete(id string) (string, error) {
	const op = "service.user.delete"

	log := u.log.With(slog.String("op", op))

	log.Info("deleting user")

	id, err := u.repository.DelById(id)
	if err != nil {
		log.Error("failed to delete user by id")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user deleted")

	return id, nil
}
