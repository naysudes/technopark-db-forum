package forum

import "github.com/naysudes/technopark-db-forum/models"

type Repository interface {
	InsertInto(f *models.Forum) error
	GetBySlug(slug string) (*models.Forum, error)
}
