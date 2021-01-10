package thread

import (
	"github.com/naysudes/technopark-db-forum/models"
)

type Usecase interface {
	CreateThread(*models.Thread) (*models.Thread, error)
	CreatePosts(string, []*models.Post) ([]*models.Post, error)
	GetBySlugOrID(string) (*models.Thread, error)
	GetPosts(string, uint64, uint64, string, bool) ([]*models.Post, error)
	Update(string, *models.Thread) (*models.Thread, error)
	Vote(slugOrId string, vote *models.Vote) (*models.Thread, error)
}
