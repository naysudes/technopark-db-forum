package user

import "github.com/naysudes/technopark-db-forum/models"

type Usecase interface {
	AddUser(string, *models.User) ([]*models.User, error)
	GetByNickname(string) (*models.User, error)
	Update(string, *models.User) error
}
