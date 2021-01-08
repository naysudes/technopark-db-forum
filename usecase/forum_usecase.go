package usecase

import (
	"github.com/naysudes/technopark-db-forum/interfaces/forum"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/user"
)

type ForumUsecase struct {
	forumRepo  forum.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
}

func NewForumUsecase(forumRepo forum.Repository, userRepo user.Repository) forum.Usecase {
	return ForumUsecase{
		forumRepo:  forumRepo,
		userRepo:   userRepo,
	}
}

func (usecase ForumUsecase) Add(forum *models.Forum) (*models.Forum, error) {
	returnForum, err := usecase.forumRepo.GetBySlug(forum.Slug)
	if err != nil && err != tools.ErrDoesntExists {
		return nil, err
	}
	if returnForum != nil {
		return returnForum, tools.ErrExistWithSlug
	}
	usr, err := usecase.userRepo.GetByNickname(forum.AdminNickname)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrUserDoesntExists
		}

		return nil, err
	}
	forum.AdminNickname = usr.Nickname
	forum.AdminID = usr.ID

	if err = usecase.forumRepo.InsertInto(forum); err != nil {
		return nil, err
	}
	return forum, nil
}

func (usecase ForumUsecase) GetBySlug(slug string) (*models.Forum, error) {
	forumBySlug, err := usecase.forumRepo.GetBySlug(slug)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}
		return nil, err
	}
	return forumBySlug, nil
}

func (usecase ForumUsecase) GetThreads(
	slug string, limit uint64, since string, desc bool) ([]*models.Thread, error) {
	if _, err := usecase.forumRepo.GetBySlug(slug); err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}
		return nil, err
	}
	threadsByForum, err := usecase.threadRepo.GetByForumSlug(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}
	return threadsByForum, nil
}

func (usecase ForumUsecase) GetUsers(
	slug string, limit uint64, since string, desc bool) ([]*models.User, error) {
	f, err := usecase.forumRepo.GetBySlug(slug)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}
	}
	usersByForum, err := usecase.userRepo.GetUsersByForum(f.ID, limit, since, desc)
	if err != nil {
		return nil, err
	}
	return usersByForum, nil
}
