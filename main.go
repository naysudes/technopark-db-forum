package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/naysudes/technopark-db-forum/database"
	"github.com/sirupsen/logrus"
	user_delivery "github.com/naysudes/technopark-db-forum/delivery/delivery"
	user_repository "github.com/naysudes/technopark-db-forum/repository/repository"
	user_usecase "github.com/naysudes/technopark-db-forum/usecase/usecase"
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

	e := echo.New()

	ur := user_repository.NewUserRepository(dbConn)

	uUC := user_usecase.NewUserUsecase(ur)
	fUC := forum_usecase.NewForumUsecase(fr, ur, pr, tr)
	tUC := thread_usecase.NewThreadUsecase(tr, ur, fr, pr, vr)
	sUC := service_usecase.NewServiceUsecase(sr)
	pUC := post_usecase.NewPostUsecase(pr, fr, vr, tr, ur)

	_ = user_delivery.NewUserHandler(e, uUC)

	e.Start(":5000")
}
