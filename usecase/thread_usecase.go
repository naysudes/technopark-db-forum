package usecase

import (
	"github.com/naysudes/technopark-db-forum/interfaces/forum"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/user"
	"strconv"
)

type ThreadUsecase struct {
	forumRepo  forum.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
}

func NewThreadUsecase(tr thread.Repository, ur user.Repository, fr forum.Repository) thread.Usecase {
	return ThreadUsecase {
		forumRepo:  fr,
		threadRepo: tr,
		userRepo:   ur,
	}
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
