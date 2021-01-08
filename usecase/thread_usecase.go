package usecase

import (
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/forum"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
	"github.com/naysudes/technopark-db-forum/interfaces/user"
	"github.com/naysudes/technopark-db-forum/interfaces/post"
	"strconv"
)

type ThreadUsecase struct {
	forumRepo  forum.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
	postRepo post.Repository
}

func NewThreadUsecase(tr thread.Repository, ur user.Repository, fr forum.Repository, pr post.Repository) thread.Usecase {
	return ThreadUsecase {
		forumRepo:  fr,
		threadRepo: tr,
		userRepo:   ur,
		postRepo: pr,
	}
}

func (tUC ThreadUsecase) CreatePosts(slugOrID string, posts []*models.Post) ([]*models.Post, error) {
	t := &models.Thread{}
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		t, err = tUC.threadRepo.GetBySlug(slugOrID)
	} else {
		t, err = tUC.threadRepo.GetByID(id)
	}
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}
		return nil, err
	}

	_, err = tUC.userRepo.CheckNicknames(posts)
	if err != nil {
		return nil, err
	}

	_, err = tUC.postRepo.CheckParentPosts(posts, t.ID)
	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return []*models.Post{}, nil
	}

	for _, p := range posts {
		p.ThreadID = t.ID
		p.Forum = t.Forum
		p.ForumID = t.ForumID
	}
	if err = tUC.postRepo.Insert(posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (tUC ThreadUsecase) GetBySlugOrID(slugOrID string) (*models.Thread, error) {
	t := &models.Thread{}

	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		t, err = tUC.threadRepo.GetBySlug(slugOrID)
	} else {
		t, err = tUC.threadRepo.GetByID(id)
	}

	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}

		return nil, err
	}
	return t, nil
}

func (tUC ThreadUsecase) CreateThread (thread *models.Thread) (*models.Thread, error) {
	forum, err := tUC.forumRepo.GetBySlug(thread.Forum)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}
		return nil, err
	}
	userAdmin, err := tUC.userRepo.GetByNickname(thread.Author)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrUserDoesntExists
		}
	}
	if thread.Slug != "" {
		returnThread, err := tUC.threadRepo.GetBySlug(thread.Slug)
		if err != nil && err != tools.ErrDoesntExists {
			return nil, err
		}
		if returnThread != nil {
			return returnThread, tools.ErrExistWithSlug
		}
	}
	thread.Forum = forum.Slug
	thread.Author = userAdmin.Nickname
	thread.AuthorID = userAdmin.ID
	thread.ForumID = forum.ID
	if err := tUC.threadRepo.InsertThread(thread); err != nil {
		return nil, err
	}
	return thread, nil
}
