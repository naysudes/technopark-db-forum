package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/naysudes/technopark-db-forum/database"
	"github.com/sirupsen/logrus"
	delivery "github.com/naysudes/technopark-db-forum/delivery"
	repository "github.com/naysudes/technopark-db-forum/repository"
	usecase "github.com/naysudes/technopark-db-forum/usecase"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
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
		logrus.Fatal(err)
		return
	}

	if err = database.InitDB(dbConn); err != nil {
		logrus.Fatal(fmt.Errorf("database init err %s", err))
		return
	}

	server := echo.New()
	server.Validator = &CustomValidator{validator: validator.New()}

	ur := repository.NewUserRepository(dbConn)
	thr := repository.NewThreadRepository(dbConn)
	fr := repository.NewForumRepository(dbConn)
	pr := repository.NewPostRepository(dbConn)
	
	uUC := usecase.NewUserUsecase(ur)
	fUC := usecase.NewForumUsecase(fr, ur)
	thUC := usecase.NewThreadUsecase(thr, ur, fr, pr)

	_ = delivery.NewThreadDelivery(server, fUC, thUC)
	_ = delivery.NewUserHandler(server, uUC)
	_ = delivery.NewForumHandler(server, thUC, fUC)


	server.Start(":5000")
}
