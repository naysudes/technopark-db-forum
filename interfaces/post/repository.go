package post

import "github.com/naysudes/technopark-db-forum/models"

type Repository interface {
	Insert([]*models.Post) error
	GetByThread(uint64, uint64, uint64, string, bool) ([]*models.Post, error)
	CheckParentPosts([]*models.Post, uint64) (bool, error)
	GetByID(uint64) (*models.Post, error)
	GetCountByForumID(uint64) (uint64, error)
	Update(*models.Post) error
}
