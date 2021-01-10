package usecase

import (
	"github.com/naysudes/technopark-db-forum/interfaces/service"
	"github.com/naysudes/technopark-db-forum/models"
)

type ServiceUsecase struct {
	serviceRepo service.Repository
}

func NewServiceUsecase(sr service.Repository) service.Usecase {
	return &ServiceUsecase{
		serviceRepo: sr,
	}
}

func (usecase *ServiceUsecase) GetStatus() (*models.Status, error) {
	forum, err := usecase.serviceRepo.GetForumCount()
	post, err := usecase.serviceRepo.GetPostCount()
	thread, err := usecase.serviceRepo.GetThreadCount()
	user, err := usecase.serviceRepo.GetUserCount()
	if err != nil {
		return nil, err
	}
	status := &models.Status{
		ForumsCount:  forum,
		PostsCount:   post,
		ThreadsCount: thread,
		UsersCount:   user,
	}
	return status, nil
}

func (usecase *ServiceUsecase) DeleteAll() error {
	err := usecase.serviceRepo.DeleteVotes()
	if err != nil {
		return err
	}
	err = usecase.serviceRepo.DeletePosts()
	if err != nil {
		return err
	}
	err = usecase.serviceRepo.DeleteThreads()
	if err != nil {
		return err
	}
	err = usecase.serviceRepo.DeleteForums()
	if err != nil {
		return err
	}
	err = usecase.serviceRepo.DeleteUsers()
	if err != nil {
		return err
	}
	return nil
}
