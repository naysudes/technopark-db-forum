package repository

import (
	"github.com/jackc/pgx"
	"github.com/naysudes/technopark-db-forum/interfaces/service"
)

type ServiceRepository struct {
	database *pgx.ConnPool
}

func NewServiceRepository(database *pgx.ConnPool) service.Repository {
	return &ServiceRepository{
		database: database,
	}
}

func (repo *ServiceRepository) GetForumCount() (uint64, error) {
	var count uint64
	if err := repo.database.QueryRow("SELECT count(*) from forums").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *ServiceRepository) GetPostCount() (uint64, error) {
	var count uint64
	if err := repo.database.QueryRow("SELECT count(*) from posts").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *ServiceRepository) GetThreadCount() (uint64, error) {
	var count uint64
	if err := repo.database.QueryRow("SELECT count(*) from threads").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *ServiceRepository) GetUserCount() (uint64, error) {
	var count uint64
	if err := repo.database.QueryRow("SELECT count(*) from users").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *ServiceRepository) DeleteUsers() error {
	if _, err := repo.database.Exec("TRUNCATE TABLE users CASCADE"); err != nil {
		return err
	}
	return nil
}

func (repo *ServiceRepository) DeleteForums() error {
	if _, err := repo.database.Exec("TRUNCATE TABLE forums CASCADE"); err != nil {
		return err
	}
	return nil
}

func (repo *ServiceRepository) DeleteThreads() error {
	if _, err := repo.database.Exec("TRUNCATE TABLE threads CASCADE"); err != nil {
		return err
	}
	return nil
}

func (repo *ServiceRepository) DeletePosts() error {
	if _, err := repo.database.Exec("TRUNCATE TABLE posts CASCADE"); err != nil {
		return err
	}
	return nil
}

func (repo *ServiceRepository) DeleteVotes() error {
	if _, err := repo.database.Exec("TRUNCATE TABLE votes CASCADE"); err != nil {
		return err
	}
	return nil
}
