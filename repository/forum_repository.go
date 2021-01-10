package repository

import (
	"github.com/jackc/pgx"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/forum"
)

type ForumRepository struct {
	database *pgx.ConnPool
}

func NewForumRepository(database *pgx.ConnPool) forum.Repository {
	return ForumRepository{
		database: database,
	}
}

func (repo ForumRepository) Insert(forum *models.Forum) error {
	if _, err := repo.database.Exec("INSERT INTO forums (slug, admin, title) VALUES ($1, $2, $3)",
		forum.Slug, forum.AdminID, forum.Title); err != nil {
		return err
	}
	return nil
}

func (repo ForumRepository) GetBySlug(slug string) (*models.Forum, error) {
	forum := &models.Forum{}
	if err := repo.database.QueryRow("SELECT f.id, f.slug, u.nickname, f.title, f.threads, f.posts FROM forums as f " +
		"JOIN users as u ON (u.id = f.admin) WHERE lower(slug) = lower($1)",
			slug).Scan(&forum.ID, &forum.Slug, &forum.AdminNickname,
		&forum.Title, &forum.ThreadsCount, &forum.PostsCount); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}
	return forum, nil
}
