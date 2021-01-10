package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/naysudes/technopark-db-forum/database"
	delivery "github.com/naysudes/technopark-db-forum/delivery"
	repository "github.com/naysudes/technopark-db-forum/repository"
	usecase "github.com/naysudes/technopark-db-forum/usecase"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func main() {
	dbConn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "forum",
			User:     "postgres",
			Password: "qweqwe",
		},
		MaxConnections: 10000,
	})
	if err != nil {
		return
	}
	if err = database.InitDB(dbConn); err != nil {
		return
	}

	server := echo.New()
	server.Validator = &Validator{validator: validator.New()}

	userRepo := repository.NewUserRepository(dbConn)
	threadRepo := repository.NewThreadRepository(dbConn)
	forumRepo := repository.NewForumRepository(dbConn)
	postRepo := repository.NewPostRepository(dbConn)
	voteRepo := repository.NewVoteRepository(dbConn)
	serviceRepo := repository.NewServiceRepository(dbConn)

	
	userUC := usecase.NewUserUsecase(userRepo)
	forumUC := usecase.NewForumUsecase(forumRepo, userRepo, threadRepo, postRepo)
	threadUC := usecase.NewThreadUsecase(threadRepo, userRepo, forumRepo, postRepo, voteRepo)
	postUC := usecase.NewPostUsecase(threadRepo, userRepo, forumRepo, postRepo)
	serviceUC := usecase.NewServiceUsecase(serviceRepo)


	_ = delivery.NewThreadDelivery(server, forumUC, threadUC)
	_ = delivery.NewUserHandler(server, userUC)
	_ = delivery.NewForumHandler(server, threadUC, forumUC)
	_ = delivery.NewPostHandler(server, postUC)
	_ = delivery.NewServiceHandler(server, serviceUC)

	server.Start(":5000")
}
