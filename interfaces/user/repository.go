package user

import "github.com/naysudes/technopark-db-forum/models"


type Repository interface {
	InsertInto(*models.User) error
	GetByNickname(string) (*models.User, error)
	CheckNicknames([]*models.Post) (bool, error)
	GetByEmail(string) (*models.User, error)
	GetUsersByForum(uint64, uint64, string, bool) ([]*models.User, error)
	Update(*models.User) error
}