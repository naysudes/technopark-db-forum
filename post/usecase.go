package post

import "github.com/naysudes/technopark-db-forum/models"

type Usecase interface {
	GetPostDetails(uint64, []string) (*models.PostFull, error)
	UpdatePost(*models.Post) (*models.Post, error)
}
