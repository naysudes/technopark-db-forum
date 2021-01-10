package usecase

import (
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/forum"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
	"github.com/naysudes/technopark-db-forum/interfaces/user"
	"github.com/naysudes/technopark-db-forum/interfaces/post"
	"github.com/naysudes/technopark-db-forum/interfaces/vote"
	"strconv"
)

type ThreadUsecase struct {
	forumRepo  forum.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
	postRepo post.Repository
	voteRepo vote.Repository
}

func NewThreadUsecase(tr thread.Repository, ur user.Repository, fr forum.Repository, pr post.Repository, vr vote.Repository) thread.Usecase {
	return &ThreadUsecase {
		forumRepo:  fr,
		threadRepo: tr,
		userRepo:   ur,
		postRepo: pr,
		voteRepo: vr,
	}
}

func (usecase *ThreadUsecase) CreatePosts(slugOrID string, posts []*models.Post) ([]*models.Post, error) {
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

func (usecase *ThreadUsecase) GetBySlugOrID(slugOrID string) (*models.Thread, error) {
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

func (usecase *ThreadUsecase) CreateThread(thread *models.Thread) (*models.Thread, error) {
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
	if err := usecase.threadRepo.Insert(thread); err != nil {
		return nil, err
	}
	return thread, nil
}

func (usecase *ThreadUsecase) GetPosts(slugOrId string, limit uint64, since uint64, sort string, descr bool) ([]*models.Post, error) {
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

func (usecase *ThreadUsecase) Update(slugOrID string, thread *models.Thread) (*models.Thread, error) {
		thread1 := &models.Thread{}

	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		thread1, err = usecase.threadRepo.GetBySlug(slugOrID)
	} else {
		thread1, err = usecase.threadRepo.GetByID(id)
	}
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}
		return nil, err
	}
	if thread.About != "" {
		thread1.About = thread.About
	}
	if thread.Title != "" {
		thread1.Title = thread.Title
	}
	if err = usecase.threadRepo.Update(thread1); err != nil {
		return nil, err
	}
	thread1.Votes, err = usecase.voteRepo.GetVotes(thread1.ID)
	if err != nil {
		return nil, err
	}
	return thread1, nil

}

func (usecase *ThreadUsecase) Vote(slugOrId string, vote *models.Vote) (*models.Thread, error) {
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
	usr, err := usecase.userRepo.GetByNickname(vote.Nickname)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrUserDoesntExists
		}
	}
	vote.ThreadID = thread.ID
	vote.UserID = usr.ID
	if err = usecase.voteRepo.Insert(vote); err != nil {
		return nil, err
	}
	thread.Votes, err = usecase.voteRepo.GetVotes(thread.ID)
	if err != nil {
		return nil, err
	}
	return thread, nil
}
