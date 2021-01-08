package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/naysudes/technopark-db-forum/models"
	// "github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/post"
)

type PostRepository struct {
	db *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) post.Repository {
	return  PostRepository{
		db: db,
	}
}

func (repo PostRepository) Insert(posts []*models.Post) error {
		queryRow := "INSERT INTO posts (author, forum, message, parent, thread) VALUES "
	for _, p := range posts {
		queryRow += fmt.Sprintf("(%d, %d, '%s', %d, %d),",
			p.AuthorID, p.ForumID, p.Message, p.ParentID, p.ThreadID)
		_, err := repo.db.Exec("INSERT INTO forums_users (user_id, forum_id) VALUES "+
			"($1, $2) ON CONFLICT DO NOTHING", p.AuthorID, p.ForumID)
		if err != nil {
			return err
		}
	}
	queryRow = queryRow[0 : len(queryRow)-1]
	queryRow += " RETURNING id, created"

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	rows, err := tx.Query(queryRow)
	if err != nil {
		return err
	}
	defer rows.Close()
	postIndex := 0
	for rows.Next() {
		if err := rows.Scan(&posts[postIndex].ID, &posts[postIndex].CreationDate); err != nil {
			return err
		}
		postIndex++
	}
	return tx.Commit()
}

func (repo PostRepository) GetByThread(uint64, uint64, uint64, string, bool) ([]*models.Post, error) {
	return nil, nil
}
func (repo PostRepository) CheckParentPosts([]*models.Post, uint64) (bool, error) {
	return true, nil
}
func (repo PostRepository) GetByID(uint64) (*models.Post, error) {
	return nil, nil
}
func (repo PostRepository) GetCountByForumID(uint64) (uint64, error) {
	return 1, nil
}
func (repo PostRepository) Update(*models.Post) error {
	return nil
}
