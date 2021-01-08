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
type UserHandler struct {
	userUsecase user.Usecase
}

func NewUserHandler(e *echo.Echo, userCase user.Usecase) *UserHandler {
	handler := &UserHandler{
		userUsecase: userCase,
	}
	e.POST("/api/user/:nickname/create", handler.CreateUser())
	e.GET("/api/user/:nickname/profile", handler.GetProfile())
	e.POST("/api/user/:nickname/profile", handler.UpdateProfile())
	return handler
}

func (handler *UserHandler) CreateUser() echo.HandlerFunc {
	type createUserRequset struct {
		Email    string `json:"email" binding:"required" validate:"email"`
		Fullname string `json:"fullname" binding:"required"`
		About    string `json:"about"`
	}

	return func(contex echo.Context) error {
		req := &createUserRequset{}
		if err := contex.Bind(req); err != nil {
			logrus.Error(fmt.Errorf("Binding error %s", err))
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{ Message: err.Error() })
		}
		if err := contex.Validate(req); err != nil {
			logrus.Error(fmt.Errorf("Validate error %s", err))
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{ Message: err.Error() })
		}
		nickname := contex.Param("nickname")
		user := &models.User{
			Email:    req.Email,
			Fullname: req.Fullname,
			About:    req.About,
		}
		returnUsers, err := handler.userUsecase.Add(nickname, user)
		if err != nil {
			if err == tools.ErrUserExistWith {
				return contex.JSON(http.StatusConflict, returnUsers)
			}
			logrus.Error(fmt.Errorf("Request error %s", err))
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return contex.JSON(http.StatusCreated, returnUsers[0])
	}
}

func (handler *UserHandler) GetProfile() echo.HandlerFunc {
	return func(contex echo.Context) error {
		nickname := contex.Param("nickname")
		returnUser, err := handler.userUsecase.GetByNickname(nickname)
		if err != nil && err != tools.ErrDoesntExists {
			logrus.Error(fmt.Errorf("Request error %s", err))
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{ Message: err.Error() })
		}
		if err == tools.ErrDoesntExists {
			return contex.JSON(http.StatusNotFound, tools.ErrorResponce{ Message: err.Error() })
		}
		return contex.JSON(http.StatusOK, returnUser)
	}
}

func (handler *UserHandler) UpdateProfile() echo.HandlerFunc {
	type updateUserRequset struct {
		Email    string `json:"email" binding:"required"`
		Fullname string `json:"fullname" binding:"required"`
		About    string `json:"about"`
	}
	return func(contex echo.Context) error {
		req := &updateUserRequset{}
		if err := contex.Bind(req); err != nil {
			logrus.Error(fmt.Errorf("Binding error %s", err))
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{ Message: err.Error() })
		}
		nickname := contex.Param("nickname")
		user := &models.User{
			Email:    req.Email,
			Fullname: req.Fullname,
			About:    req.About,
		}
		err := handler.userUsecase.Update(nickname, user)
		if err != nil {
			if err == tools.ErrUserExistWith {
				return contex.JSON(http.StatusConflict, tools.ErrorResponce{ Message: err.Error() })
			}
			if err == tools.ErrUserDoesntExists {
				return contex.JSON(http.StatusNotFound, tools.ErrorResponce{ Message: err.Error() })
			}
			logrus.Error(fmt.Errorf("Request error %s", err))
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{ Message: err.Error() })
		}
		return contex.JSON(http.StatusOK, user)
	}
}
