package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
)

type ThreadRepository struct {
	database *pgx.ConnPool
}

func NewThreadRepository(database *pgx.ConnPool) thread.Repository {
	return &ThreadRepository{
		database: database,
	}
}

func (repo *ThreadRepository) Insert(t *models.Thread) error {
	if err := repo.database.QueryRow("INSERT INTO threads "+
		"(slug, author, title, message, forum, created) "+
		"VALUES (NULLIF ($1, ''), $2, $3, $4, $5, $6) RETURNING id",
		t.Slug, t.AuthorID, t.Title, t.About, t.ForumID, t.CreationDate).
		Scan(&t.ID); err != nil {
		return err
	}

	return nil
}

func (repo *ThreadRepository) GetByID(id uint64) (*models.Thread, error) {
	thread := &models.Thread{}
	if err := repo.database.QueryRow("SELECT t.id, u.nickname, t.created, t.forum, f.slug, t.message, "+
		"coalesce (t.slug, ''), t.title, t.votes FROM threads AS t "+
		"JOIN users AS u ON (t.author = u.id) "+
		"JOIN forums AS f ON (f.id = t.forum) WHERE t.id = $1", id).
		Scan(&thread.ID, &thread.Author, &thread.CreationDate, &thread.ForumID, &thread.Forum, &thread.About, &thread.Slug,
			&thread.Title, &thread.Votes); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}
	return thread, nil
}

func (repo *ThreadRepository) GetBySlug(slug string) (*models.Thread, error) {
	thread := &models.Thread{}
	if err := repo.database.QueryRow("SELECT t.id, u.nickname, t.created, t.forum, f.slug, t.message, "+
		"coalesce (t.slug, ''), t.title, t.votes FROM threads AS t "+
		"JOIN users AS u ON (t.author = u.id) "+
		"JOIN forums AS f ON (f.id = t.forum) WHERE lower(t.slug) = lower($1)", slug).
		Scan(&thread.ID, &thread.Author, &thread.CreationDate, &thread.ForumID, &thread.Forum, &thread.About, &thread.Slug,
			&thread.Title, &thread.Votes); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}
	return thread, nil
}

func (repo *ThreadRepository) GetByForumSlug(
	slug string, limit uint64, since string, desc bool) ([]*models.Thread, error) {
	threadsByForum := []*models.Thread{}
	queryStr := "SELECT t.id, u.nickname, t.created, f.slug, t.message, " +
		"coalesce (t.slug, ''), t.title, t.votes FROM threads AS t " +
		"JOIN users AS u ON (t.author = u.id) " +
		"JOIN forums AS f ON (f.id = t.forum) WHERE lower(f.slug) = lower($1)"

	orderStr := " ORDER BY t.created"
	if desc {
		orderStr += " DESC"
	}
	if limit != 0 {
		orderStr += fmt.Sprintf(" LIMIT %d", limit)
	}
	var rows *pgx.Rows
	var err error
	if since != "" {
		if desc {
			queryStr += " AND t.created <= $2"
		} else {
			queryStr += " AND t.created >= $2"
		}
		rows, err = repo.database.Query(queryStr + orderStr, slug, since)
	} else {
		rows, err = repo.database.Query(queryStr + orderStr, slug)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		th := &models.Thread{}
		err = rows.Scan(&th.ID, &th.Author, &th.CreationDate, &th.Forum, &th.About, &th.Slug, &th.Title, &th.Votes)
		if err != nil {
			return nil, err
		}
		threadsByForum = append(threadsByForum, th)
	}

	return threadsByForum, nil
}

func (repo *ThreadRepository) GetCountByForumID(id uint64) (uint64, error) {
	var count uint64
	if err := repo.database.QueryRow("SELECT count(*) FROM threads WHERE forum = $1 ", id).
		Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *ThreadRepository) Update(thread *models.Thread) error {
	if _, err := repo.database.Exec("UPDATE threads SET message = $2, title = $3 WHERE id = $1",
		thread.ID, thread.About, thread.Title); err != nil {
		return err
	}
	return nil
}
