package repository

import (
	"github.com/jackc/pgx"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/interfaces/vote"
)

type VoteRepository struct {
	database *pgx.ConnPool
}

func NewVoteRepository(database *pgx.ConnPool) vote.Repository {
	return &VoteRepository{
		database: database,
	}
}

func (repo *VoteRepository) GetVotes(id uint64) (int64, error) {
	var votes int64
	if err := repo.database.QueryRow("SELECT votes "+
		"FROM threads WHERE id = $1", id).Scan(&votes); err != nil {
		return 0, err
	}
	return votes, nil
}

func (repo *VoteRepository) Insert(vote *models.Vote) error {
	curentId := uint64(0)
	err := repo.database.QueryRow("SELECT id FROM votes "+
		"WHERE thread = $1 and author = $2", vote.ThreadID, vote.UserID).Scan(&curentId)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}
	if err == pgx.ErrNoRows {
		if _, err := repo.database.Exec("INSERT INTO votes (author, thread, vote) "+
			"VALUES ($1, $2, $3)",
			vote.UserID, vote.ThreadID, vote.Voice); err != nil {
			return err
		}
		return nil
	}
	if _, err = repo.database.Exec("UPDATE votes SET vote = $2 "+
		"WHERE id = $1", curentId, vote.Voice); err != nil {
		return err
	}
	return nil
}
