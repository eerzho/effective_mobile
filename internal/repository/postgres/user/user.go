package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"effective_mobile/internal/conn/postgres"
	"effective_mobile/internal/domain"
	"effective_mobile/internal/exception"
)

type Repo struct {
	conn *postgres.Conn
}

func New(con *postgres.Conn) *Repo {
	return &Repo{conn: con}
}

func (r *Repo) List(page int, name, surname string) ([]*domain.User, error) {
	const op = "repository.postgres.user.List"

	offset := (page - 1) * 10
	var users []*domain.User

	query := `SELECT id, name, surname, patronymic, gender, country_id, age FROM users WHERE 1=1`
	var args []interface{}

	if name != "" {
		query += ` AND name LIKE '%' || $` + strconv.Itoa(len(args)+1) + ` || '%'`
		args = append(args, name)
	}
	if surname != "" {
		query += ` AND surname LIKE '%' || $` + strconv.Itoa(len(args)+1) + ` || '%'`
		args = append(args, surname)
	}

	query += fmt.Sprintf(" ORDER BY name LIMIT 10 OFFSET %d", offset)

	rows, err := r.conn.DB.Query(query, args...)
	if err != nil {
		return []*domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Patronymic, &user.Gender, &user.CountryId, &user.Age); err != nil {
			return []*domain.User{}, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return []*domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (r *Repo) GetById(id string) (*domain.User, error) {
	const op = "repository.postgres.user.GetById"

	var user domain.User

	query := `SELECT id, name, surname, patronymic, gender, country_id, age FROM users WHERE id = $1`
	err := r.conn.DB.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Surname, &user.Patronymic, &user.Gender, &user.CountryId, &user.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.User{}, fmt.Errorf("%s: %w", op, exception.ErrUserNotFound)
		}
		return &domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *Repo) Exists(id string) (bool, error) {
	const op = "repository.postgres.user.Exists"

	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := r.conn.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, exception.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (r *Repo) Save(name, surname string, patronymic, gender, countryId *string, age *int) (*domain.User, error) {
	const op = "repository.postgres.user.Save"

	var user domain.User

	query := `INSERT INTO users (name, surname, patronymic, gender, country_id, age) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.conn.DB.QueryRow(query, name, surname, patronymic, gender, countryId, age).Scan(&user.Id)
	if err != nil {
		return &domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user.Name = name
	user.Surname = surname
	user.Patronymic = patronymic
	user.Gender = gender
	user.CountryId = countryId
	user.Age = age

	return &user, nil
}

func (r *Repo) Update(id, name, surname string, patronymic, gender, countryId *string, age *int) (*domain.User, error) {
	const op = "repository.postgres.user.Update"

	query := `UPDATE users SET name = $1, surname = $2, patronymic = $3, gender = $4, country_id = $5, age = $6 WHERE id = $7`
	_, err := r.conn.DB.Exec(query, name, surname, patronymic, gender, countryId, age, id)
	if err != nil {
		return &domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return r.GetById(id)
}

func (r *Repo) DelById(id string) (string, error) {
	const op = "repository.postgres.user.DelById"

	query := `DELETE FROM users WHERE id = $1`
	result, err := r.conn.DB.Exec(query, id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if affected == 0 {
		return "", fmt.Errorf("%s: %w", op, exception.ErrUserNotFound)
	}

	return id, nil
}
