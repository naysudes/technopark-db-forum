package post

import "github.com/naysudes/technopark-db-forum/models"

type Usecase interface {
	GetDetails(uint64, []string) (*models.PostDetailed, error)
	Update(*models.Post) (*models.Post, error)
}
