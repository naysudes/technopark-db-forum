package usecase

import (
	"github.com/naysudes/technopark-db-forum/interfaces/forum"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/user"
	"github.com/naysudes/technopark-db-forum/interfaces/post"
)

type ForumUsecase struct {
	forumRepo  forum.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
	postRepo   post.Repository
}

func NewForumUsecase(fr forum.Repository, ur user.Repository, tr thread.Repository, pr post.Repository) forum.Usecase {
	return &ForumUsecase{
		forumRepo:  fr,
		threadRepo: tr,
		userRepo:   ur,
		postRepo:   pr,
	}
}

func (usecase *ForumUsecase) Add(forum *models.Forum) (*models.Forum, error) {
	forumBySlug, err := usecase.forumRepo.GetBySlug(forum.Slug)
	if err != nil && err != tools.ErrDoesntExists {
		return nil, err
	}
	if forumBySlug != nil {
		return forumBySlug, tools.ErrExistWithSlug
	}
	u, err := usecase.userRepo.GetByNickname(forum.AdminNickname)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrUserDoesntExists
		}
		return nil, err
	}
	forum.AdminNickname = u.Nickname
	forum.AdminID = u.ID
	if err = usecase.forumRepo.Insert(forum); err != nil {
		return nil, err
	}
	return forum, nil
}

func (usecase *ForumUsecase) GetBySlug(slug string) (*models.Forum, error) {
	forum, err := usecase.forumRepo.GetBySlug(slug)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}
		return nil, err
	}
	return forum, nil
}

func (usecase *ForumUsecase) GetUsers(
	slug string, limit uint64, since string, desc bool) ([]*models.User, error) {
	forum, err := usecase.forumRepo.GetBySlug(slug)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}
	}
	users, err := usecase.userRepo.GetUsersByForum(forum.ID, limit, since, desc)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (usecase *ForumUsecase) GetThreadsByForum(slug string, limit uint64, since string, desc bool) ([]*models.Thread, error) {
	if _, err := usecase.forumRepo.GetBySlug(slug); err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}
		return nil, err
	}
	threads, err := usecase.threadRepo.GetByForumSlug(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}
	return threads, nil
}
