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

func (usecase ThreadUsecase) CreatePosts(slugOrID string, posts []*models.Post) ([]*models.Post, error) {
	t := &models.Thread{}
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		t, err = usecase.threadRepo.GetBySlug(slugOrID)
	} else {
		t, err = usecase.threadRepo.GetByID(id)
	}
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}
		return nil, err
	}

	_, err = usecase.userRepo.CheckNicknames(posts)
	if err != nil {
		return nil, err
	}

	_, err = usecase.postRepo.CheckParentPosts(posts, t.ID)
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
	if err = usecase.postRepo.Insert(posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (usecase ThreadUsecase) GetBySlugOrID(slugOrID string) (*models.Thread, error) {
	t := &models.Thread{}

	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		t, err = usecase.threadRepo.GetBySlug(slugOrID)
	} else {
		t, err = usecase.threadRepo.GetByID(id)
	}

	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}

		return nil, err
	}
	return t, nil
}

func (usecase ThreadUsecase) CreateThread(thread *models.Thread) (*models.Thread, error) {
	forum, err := usecase.forumRepo.GetBySlug(thread.Forum)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}
		return nil, err
	}
	userAdmin, err := usecase.userRepo.GetByNickname(thread.Author)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrUserDoesntExists
		}
	}
	if thread.Slug != "" {
		returnThread, err := usecase.threadRepo.GetBySlug(thread.Slug)
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
	if err := usecase.threadRepo.InsertThread(thread); err != nil {
		return nil, err
	}
	return thread, nil
}

func (usecase ThreadUsecase) GetPosts(slugOrId string, limit uint64, since uint64, sort string, descr bool) ([]*models.Post, error) {
	thread := &models.Thread{}
	id, err := strconv.ParseUint(slugOrId, 10, 64)
	if err != nil {
		thread, err = usecase.threadRepo.GetBySlug(slugOrId)
	} else {
		thread, err = usecase.threadRepo.GetByID(id)
	}
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}
		return nil, err
	}
	postsByThread, err := usecase.postRepo.GetByThread(thread.ID, limit, since, sort, descr)
	if err != nil {
		return nil, err
	}
	return postsByThread, err
}

func (usecase ThreadUsecase) Update(slugOrID string, thread *models.Thread) (*models.Thread, error) {
		t := &models.Thread{}

	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		t, err = usecase.threadRepo.GetBySlug(slugOrID)
	} else {
		t, err = usecase.threadRepo.GetByID(id)
	}
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}
		return nil, err
	}
	if thread.About != "" {
		t.About = thread.About
	}
	if thread.Title != "" {
		t.Title = thread.Title
	}
	if err = usecase.threadRepo.Update(t); err != nil {
		return nil, err
	}
	// t.Votes, err = usecase.voteRepo.GetThreadVotes(t.ID)
	// if err != nil {
	// 	return nil, err
	// }
	return t, nil

}
