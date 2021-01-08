package repository

import (
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

func (repo PostRepository) InsertInto([]*models.Post) error {
	return nil
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
