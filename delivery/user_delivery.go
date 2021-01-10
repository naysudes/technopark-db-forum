package delivery

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/user"
	"github.com/sirupsen/logrus"
	"net/http"
)
type UserDelivery struct {
	userUsecase user.Usecase
}

func NewUserHandler(e *echo.Echo, userCase user.Usecase) *UserDelivery {
	handler := &UserDelivery{
		userUsecase: userCase,
	}
	e.POST("/api/user/:nickname/create", handler.CreateUser())
	e.GET("/api/user/:nickname/profile", handler.GetProfile())
	e.POST("/api/user/:nickname/profile", handler.UpdateProfile())
	return handler
}

func (delivery *UserDelivery) CreateUser() echo.HandlerFunc {
	type createUserRequset struct {
		Email    string `json:"email" binding:"required" validate:"email"`
		Fullname string `json:"fullname" binding:"required"`
		About    string `json:"about"`
	}

	return func(context echo.Context) error {
		req := &createUserRequset{}
		if err := context.Bind(req); err != nil {
			logrus.Error(fmt.Errorf("Binding error %s", err))
			return context.JSON(http.StatusBadRequest, tools.ErrorResponce{ Message: err.Error() })
		}
		if err := context.Validate(req); err != nil {
			logrus.Error(fmt.Errorf("Validate error %s", err))
			return context.JSON(http.StatusBadRequest, tools.ErrorResponce{ Message: err.Error() })
		}
		nickname := context.Param("nickname")
		user := &models.User{
			Email:    req.Email,
			Fullname: req.Fullname,
			About:    req.About,
		}
		returnUsers, err := delivery.userUsecase.Add(nickname, user)
		if err != nil {
			if err == tools.ErrUserExistWith {
				return context.JSON(http.StatusConflict, returnUsers)
			}
			logrus.Error(fmt.Errorf("Request error %s", err))
			return context.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.JSON(http.StatusCreated, returnUsers[0])
	}
}

func (delivery *UserDelivery) GetProfile() echo.HandlerFunc {
	return func(context echo.Context) error {
		nickname := context.Param("nickname")
		returnUser, err := delivery.userUsecase.GetByNickname(nickname)
		if err != nil && err != tools.ErrDoesntExists {
			logrus.Error(fmt.Errorf("Request error %s", err))
			return context.JSON(http.StatusBadRequest, tools.ErrorResponce{ Message: err.Error() })
		}
		if err == tools.ErrDoesntExists {
			return context.JSON(http.StatusNotFound, tools.ErrorResponce{ Message: err.Error() })
		}
		return context.JSON(http.StatusOK, returnUser)
	}
}

func (delivery *UserDelivery) UpdateProfile() echo.HandlerFunc {
	type updateUserRequset struct {
		Email    string `json:"email" binding:"required"`
		Fullname string `json:"fullname" binding:"required"`
		About    string `json:"about"`
	}
	return func(context echo.Context) error {
		req := &updateUserRequset{}
		if err := context.Bind(req); err != nil {
			logrus.Error(fmt.Errorf("Binding error %s", err))
			return context.JSON(http.StatusBadRequest, tools.ErrorResponce{ Message: err.Error() })
		}
		nickname := context.Param("nickname")
		user := &models.User{
			Email:    req.Email,
			Fullname: req.Fullname,
			About:    req.About,
		}
		err := delivery.userUsecase.Update(nickname, user)
		if err != nil {
			if err == tools.ErrUserExistWith {
				return context.JSON(http.StatusConflict, tools.ErrorResponce{ Message: err.Error() })
			}
			if err == tools.ErrUserDoesntExists {
				return context.JSON(http.StatusNotFound, tools.ErrorResponce{ Message: err.Error() })
			}
			logrus.Error(fmt.Errorf("Request error %s", err))
			return context.JSON(http.StatusBadRequest, tools.ErrorResponce{ Message: err.Error() })
		}
		return context.JSON(http.StatusOK, user)
	}
}
