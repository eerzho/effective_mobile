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

type User struct {
	conn *postgres.Conn
}

func New(con *postgres.Conn) *User {
	return &User{conn: con}
}

func (u *User) List(page int, name, surname string) ([]domain.User, error) {
	const op = "repository.postgres.user.List"

	offset := (page - 1) * 10
	var users []domain.User

	query := `SELECT id, name, surname, patronymic, sex, nationality, age FROM users WHERE 1=1`
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

	rows, err := u.conn.DB.Query(query, args...)
	if err != nil {
		return []domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Patronymic, &user.Sex, &user.Nationality, &user.Age); err != nil {
			return []domain.User{}, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return []domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (u *User) GetById(id string) (domain.User, error) {
	const op = "repository.postgres.user.GetById"

	var user domain.User

	query := `SELECT id, name, surname, patronymic, sex, nationality, age FROM users WHERE id = $1`
	err := u.conn.DB.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Surname, &user.Patronymic, &user.Sex, &user.Nationality, &user.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", op, exception.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) Exists(id string) (bool, error) {
	const op = "repository.postgres.user.Exists"

	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := u.conn.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, exception.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (u *User) Save(name, surname string, patronymic, sex, nationality *string, age *int) (domain.User, error) {
	const op = "repository.postgres.user.Save"

	var user domain.User

	query := `INSERT INTO users (name, surname, patronymic, sex, nationality, age) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := u.conn.DB.QueryRow(query, name, surname, patronymic, sex, nationality, age).Scan(&user.Id)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user.Name = name
	user.Surname = surname
	user.Patronymic = patronymic
	user.Sex = sex
	user.Nationality = nationality
	user.Age = age

	return user, nil
}

func (u *User) Update(id, name, surname string, patronymic, sex, nationality *string, age *int) (domain.User, error) {
	const op = "repository.postgres.user.Update"

	query := `UPDATE users SET name = $1, surname = $2, patronymic = $3, sex = $4, nationality = $5, age = $6 WHERE id = $7`
	_, err := u.conn.DB.Exec(query, name, surname, patronymic, sex, nationality, age, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return u.GetById(id)
}
func (u *User) DelById(id string) (string, error) {
	const op = "repository.postgres.user.DelById"

	query := `DELETE FROM users WHERE id = $1`
	result, err := u.conn.DB.Exec(query, id)
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
