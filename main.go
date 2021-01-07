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
	uUC := usecase.NewUserUsecase(ur)
	_ = delivery.NewUserHandler(server, uUC)

	fr := repository.NewForumRepository(dbConn)
	uUC := usecase.NewForumUsecase(ur)
	_ = delivery.NewForumHandler(server, uUC)

	server.Start(":5000")
}
