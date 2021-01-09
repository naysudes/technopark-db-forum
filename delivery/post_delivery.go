package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	// "github.com/naysudes/technopark-db-forum/interfaces/user"
	"github.com/naysudes/technopark-db-forum/interfaces/post"
	"net/http"
	"strconv"
	"strings"
	// "encoding/json"
)

type PostHandler struct {
	postUsecase post.Usecase
}

func NewPostHandler(e *echo.Echo, pUC post.Usecase) PostHandler {
	ph := PostHandler{
		postUsecase: pUC,
	}

	e.GET("/api/post/:id/details", ph.GetPostDetails())
	e.POST("/api/post/:id/details", ph.UpdatePost())

	return ph
}

func (ph PostHandler) GetPostDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		related := strings.Split(c.QueryParam("related"), ",")
		returnPost, err := ph.postUsecase.GetDetails(id, related)
		if err != nil {
			if err == tools.ErrPostDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, returnPost)
	}
}

func (ph PostHandler) UpdatePost() echo.HandlerFunc {
	type updatePostReq struct {
		Message string `json:"message" binding:"require"`
	}
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		req := &updatePostReq{}
		if err = c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		p := &models.Post{
			ID: id,
			Message: req.Message,
		}
		updPost, err := ph.postUsecase.Update(p)
		if err != nil {
			if err == tools.ErrPostDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, updPost)
	}
}
