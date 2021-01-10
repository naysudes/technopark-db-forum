package service

import "github.com/naysudes/technopark-db-forum/models"

type Usecase interface {
	GetStatus() (*models.Status, error)
	DeleteAll() error
}