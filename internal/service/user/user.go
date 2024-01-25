package user

import (
	"fmt"
	"log/slog"

	"effective_mobile/internal/domain"
)

type User struct {
	log        *slog.Logger
	repository Repository
}

func New(log *slog.Logger, repository Repository) *User {
	return &User{log: log, repository: repository}
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

	user, err := u.repository.Save(name, surname, patronymic, nil, nil, nil)
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
