package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/naysudes/technopark-db-forum/thread"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/forum"
	"net/http"
	"strconv"
)

type ForumHandler struct {
	forumUC  forum.Usecase
	threadUC thread.Usecase
}

func NewForumHandler(e *echo.Echo, fUC forum.Usecase) *ForumHandler {
	fh := &ForumHandler{ forumUC:  fUC }

	e.POST("/api/forum/create", fh.CreateForum())
	e.GET("/api/forum/:slug/details", fh.GetForumDetails())
	e.GET("/api/forum/:slug/users", fh.GetForumUsers())

	return fh
}

func (fh *ForumHandler) CreateForum() echo.HandlerFunc {
	type createForumRequest struct {
		Slug  string `json:"slug" binding:"required"`
		Title string `json:"title" binding:"required"`
		User  string `json:"user" binding:"required"`
	}

	return func(c echo.Context) error {
		req := &createForumRequest{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		reqForum := &models.Forum{
			Slug:          req.Slug,
			Title:         req.Title,
			AdminNickname: req.User,
		}

		returnForum, err := fh.forumUC.AddForum(reqForum)
		if err != nil {
			if err == tools.ErrExistWithSlug {
				return c.JSON(http.StatusConflict, returnForum)
			}
			if err == tools.ErrUserDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, returnForum)
	}
}

func (fh *ForumHandler) GetForumDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug")

		returnForum, err := fh.forumUC.GetForumBySlug(slug)
		if err != nil {
			if err == tools.ErrForumDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, returnForum)
	}
}

func (fh *ForumHandler) GetForumUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug")

		limit := uint64(0)
		var err error
		if l := c.QueryParam("limit"); l != "" {
			limit, err = strconv.ParseUint(l, 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
		}
		since := c.QueryParam("since")

		desc := false
		if descVal := c.QueryParam("desc"); descVal == "true" {
			desc = true
		}

		returnUsers, err := fh.forumUC.GetForumUsers(slug, limit, since, desc)
		if err != nil {
			if err == tools.ErrForumDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, returnUsers)
	}
}
