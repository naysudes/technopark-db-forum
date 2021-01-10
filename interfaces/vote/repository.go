package vote

import "github.com/naysudes/technopark-db-forum/models"

type Repository interface {
	GetVotes(id uint64) (int64, error)
	Insert(*models.Vote) error
}