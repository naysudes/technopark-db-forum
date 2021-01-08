package forum

import "github.com/naysudes/technopark-db-forum/models"

type Usecase interface {
	Add(*models.Forum) (*models.Forum, error)
	GetBySlug(string) (*models.Forum, error)
	GetUsers(string, uint64, string, bool) ([]*models.User, error)
	GetThreads(string, uint64, string, bool) ([]*models.Thread, error)
}
