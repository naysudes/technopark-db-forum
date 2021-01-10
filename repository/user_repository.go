package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/user"
	"strings"
)

type UserRepository struct {
	database *pgx.ConnPool
}

func NewUserRepository(database *pgx.ConnPool) user.Repository {
	return &UserRepository{ database: database }
}

func (repo *UserRepository) Insert(usr *models.User) error {
	if _, err := repo.database.Exec("INSERT INTO users (nickname, email, fullname, about) VALUES ($1, $2, $3, $4)",
	usr.Nickname, usr.Email, usr.Fullname, usr.About); err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetUsersByForum(
	id uint64, limit uint64, since string, desc bool) ([]*models.User, error) {
	usrs := []*models.User{}
	queryString := "SELECT u.nickname, u.email, u.fullname, u.about FROM forums_users fu " +
		"JOIN users u ON (fu.user_id = u.id) " +
		"WHERE fu.forum_id = $1"
	orderbyString := " ORDER BY lower(u.nickname)"
	if desc {
		orderbyString += " DESC"
	}
	if limit != 0 {
		orderbyString += fmt.Sprintf(" LIMIT %d", limit)
	}
	var rows *pgx.Rows
	var err error
	if since != "" {
		if desc {
			queryString += " AND lower(u.nickname) < lower($2)"
		} else {
			queryString += " AND lower(u.nickname) > lower($2)"
		}
		rows, err = repo.database.Query(queryString + orderbyString, id, since)
	} else {
		rows, err = repo.database.Query(queryString + orderbyString, id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		usr := &models.User{}
		err = rows.Scan(&usr.Nickname, &usr.Email, &usr.Fullname, &usr.About)
		if err != nil {
			return nil, err
		}
		usrs = append(usrs, usr)
	}
	return usrs, nil
}

func (repo *UserRepository) GetByNickname(nickname string) (*models.User, error) {
	usr := &models.User{}
	if err := repo.database.QueryRow("SELECT id, nickname, email, fullname, about FROM users " +
		"WHERE lower(nickname) = lower($1)", nickname).Scan(&usr.ID, &usr.Nickname, &usr.Email, &usr.Fullname,
		&usr.About); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}
	return usr, nil
}

func (repo *UserRepository) Update(usr *models.User) error {
	if _, err := repo.database.Exec("UPDATE users SET email = $2, fullname = $3, about = $4 WHERE lower(nickname) = lower($1)",
		usr.Nickname, usr.Email, usr.Fullname, usr.About); err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
	usr := &models.User{}
	if err := repo.database.QueryRow("SELECT id, nickname, email, fullname, about FROM users WHERE lower(email) = lower($1)",
		email).Scan(&usr.ID, &usr.Nickname, &usr.Email, &usr.Fullname, &usr.About); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}
	return usr, nil
}

func (repo *UserRepository) CheckNicknames(posts []*models.Post) (bool, error) {
	rows, err := repo.database.Query("SELECT id, lower(nickname) FROM users")
	if err != nil {
		return false, err
	}
	defer rows.Close()
	nicknames := make(map[string]uint64)
	for rows.Next() {
		n := ""
		var id uint64
		if err := rows.Scan(&id, &n); err != nil {
			return false, err
		}
		nicknames[n] = id
	}

	for _, p := range posts {
		id := nicknames[strings.ToLower(p.Author)]
		if id == 0 {
			return false, tools.ErrUserDoesntExists
		}
		p.AuthorID = id
	}
	return true, nil
}
