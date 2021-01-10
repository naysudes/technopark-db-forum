package usecase

import (
	"github.com/naysudes/technopark-db-forum/tools"
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
func (usecase PostUsecase) GetDetails(id uint64, related []string) (*models.PostDetailed, error) {
	post, err := usecase.postRepo.GetByID(id)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrPostDoesntExists
		}
		return nil, err
	}
	postDetails := &models.PostDetailed{Post: post}

	for _, rel := range related {
		if (rel == "user") {
			usr, err := usecase.userRepo.GetByNickname(post.Author)
			if err != nil {
				if err == tools.ErrDoesntExists {
					return nil, tools.ErrUserDoesntExists
				}
				return nil, err
			}
			postDetails.Author = usr
		}
		if (rel == "thread") {
			thead, err := usecase.threadRepo.GetByID(post.ThreadID)
			if err != nil {
				if err == tools.ErrDoesntExists {
					return nil, tools.ErrThreadDoesntExists
				}
				return nil, err
			}
			postDetails.Thread = thead
		}
		if (rel == "forum") {
			forum, err := usecase.forumRepo.GetBySlug(post.Forum)
			if err != nil {
				if err == tools.ErrDoesntExists {
					return nil, tools.ErrForumDoesntExists
				}
				return nil, err
			}
			postDetails.Forum = forum
		}
	}
	return postDetails, nil


}
func (usecase PostUsecase) Update(post *models.Post) (*models.Post, error) {
	updated, err := usecase.postRepo.GetByID(post.ID)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrPostDoesntExists
		}

		return nil, err
	}
	if post.Message == "" || post.Message == updated.Message {
		return updated, nil
	}
	updated.Message = post.Message
	updated.IsEdited = true

	if err = usecase.postRepo.Update(updated); err != nil {
		return nil, err
	}
	return updated, nil
}