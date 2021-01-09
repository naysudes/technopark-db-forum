package thread

import (
	"github.com/naysudes/technopark-db-forum/models"
)

type Repository interface {
	InsertThread(*models.Thread) error
	GetByID(uint64) (*models.Thread, error)
	GetBySlug(string) (*models.Thread, error)
	GetByForumSlug(string, uint64, string, bool) ([]*models.Thread, error)
	// GetCountByForumID(uint64) (uint64, error)
	Update(*models.Thread) error
}
