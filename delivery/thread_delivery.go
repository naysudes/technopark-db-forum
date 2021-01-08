package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/forum"
	// "github.com/naysudes/technopark-db-forum/interfaces/user"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
	"net/http"
	// "strconv"
	// "encoding/json"
)

type ThreadDelivery struct {
	forumUseCase  forum.Usecase
	threadUseCase thread.Usecase
}

func NewThreadDelivery(e *echo.Echo, forumUC forum.Usecase, threadUC thread.Usecase) ThreadDelivery {
	handler := ThreadDelivery{ forumUseCase:  forumUC, threadUseCase: threadUC }
	e.POST("/api/thread/:slug_or_id/create", handler.CreatePosts())
	e.GET("/api/thread/:slug_or_id/details", handler.GetDetails())
	// e.POST("/api/thread/:slug_or_id/details", handler.Update())
	// e.GET("/api/thread/:slug_or_id/posts", handler.GetPosts())
	// e.POST("/api/thread/:slug_or_id/vote", handler.Vote())
	return handler
}

func (th ThreadDelivery) CreatePosts() echo.HandlerFunc {
	type createReq struct {
		Author  string `json:"author" binding:"required"`
		Message string `json:"message" binding:"required"`
		Parent  uint64 `json:"parent" binding:"required"`
	}
  	return func(contex echo.Context) error {
		req := []createReq{}
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		slugOrID := contex.Param("slug_or_id")

		posts := make([]models.Post, 0, len(req))
		for _, r := range req {
			post := models.Post{
				Author:   r.Author,
				Message:  r.Message,
				ParentID: r.Parent,
				IsEdited: false,
			}
			posts = append(posts, post)
		}
		returnPosts, err := th.threadUseCase.CreatePosts(slugOrID, posts)
		if err != nil {
			if err == tools.ErrThreadDoesntExists || err == tools.ErrUserDoesntExists {
				return contex.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			return contex.JSON(http.StatusConflict, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return contex.JSON(http.StatusCreated, returnPosts)
	}
}

func (th ThreadDelivery) GetDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		slugOrID := c.Param("slug_or_id")
		threadBySlugorId, err := th.threadUseCase.GetBySlugOrID(slugOrID)
		if err != nil {
			if err == tools.ErrThreadDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, threadBySlugorId)
	}
}

func (th ThreadDelivery) Update() {

}

func (th ThreadDelivery) GetPosts() {

}

func (th ThreadDelivery) Vote() {

}
