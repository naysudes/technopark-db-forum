package usecase

import (
	// "github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/interfaces/forum"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
	"github.com/naysudes/technopark-db-forum/interfaces/user"
	"github.com/naysudes/technopark-db-forum/interfaces/post"
)

type PostUsecase struct {
	forumRepo  forum.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
	postRepo post.Repository
}

func NewPostUsecase(tr thread.Repository, ur user.Repository, fr forum.Repository, pr post.Repository) post.Usecase {
	return PostUsecase {
		forumRepo:  fr,
		threadRepo: tr,
		userRepo:   ur,
		postRepo: pr,
	}
}
func (pUC PostUsecase) GetDetails(uint64, []string) (*models.PostFull, error) {
	return nil, nil

}
func (pUC PostUsecase) Update(*models.Post) (*models.Post, error) {
	return nil, nil
}