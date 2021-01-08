package repository

import (
	"github.com/jackc/pgx"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
)

type ThreadRepository struct {
	db *pgx.ConnPool
}

func NewThreadRepository(db *pgx.ConnPool) thread.Repository {
	return ThreadRepository{
		db: db,
	}
}

func (tr ThreadRepository) GetBySlug(slug string) (*models.Thread, error) {
	t := &models.Thread{}
	if err := tr.db.QueryRow("SELECT t.id, u.nickname, t.created, t.forum, f.slug, t.message, "+
		"coalesce (t.slug, ''), t.title, t.votes FROM threads AS t "+
		"JOIN users AS u ON (t.author = u.id) "+
		"JOIN forums AS f ON (f.id = t.forum) WHERE lower(t.slug) = lower($1)", slug).
		Scan(&t.ID, &t.Author, &t.CreationDate, &t.ForumID, &t.Forum, &t.About, &t.Slug,
			&t.Title, &t.Votes); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}
	return t, nil
}

func (tr ThreadRepository) GetByID(id uint64) (*models.Thread, error) {
	t := &models.Thread{}
	if err := tr.db.QueryRow("SELECT t.id, u.nickname, t.created, t.forum, f.slug, t.message, "+
		"coalesce (t.slug, ''), t.title, t.votes FROM threads AS t "+
		"JOIN users AS u ON (t.author = u.id) "+
		"JOIN forums AS f ON (f.id = t.forum) WHERE t.id = $1", id).
		Scan(&t.ID, &t.Author, &t.CreationDate, &t.ForumID, &t.Forum, &t.About, &t.Slug,
			&t.Title, &t.Votes); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}
	return t, nil
}

func (tr ThreadRepository) GetByForumSlug(string, uint64, string, bool) ([]*models.Thread, error) {
	return nil, nil
}

func (tr ThreadRepository) InsertThread(t *models.Thread) error {
	if err := tr.db.QueryRow("INSERT INTO threads (slug, author, title, message, forum, created) " +
		"VALUES (NULLIF ($1, ''), $2, $3, $4, $5, $6) RETURNING id",
		t.Slug, t.AuthorID, t.Title, t.About, t.ForumID, t.CreationDate).
		Scan(&t.ID); err != nil {
		return err
	}
	return nil
}