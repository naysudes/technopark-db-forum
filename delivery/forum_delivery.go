package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/forum"
	"net/http"
	"strconv"
	"time"
)

type ForumDelivery struct {
	forumUseCase  forum.Usecase
	threadUseCase thread.Usecase
}

func NewForumHandler(e *echo.Echo, thUC thread.Usecase, fUC forum.Usecase) *ForumDelivery {
	handler := &ForumDelivery{ forumUseCase:  fUC, threadUseCase: thUC}

	e.POST("/api/forum/create", handler.CreateForum())
	e.GET("/api/forum/:slug/details", handler.GetForumDetails())
	e.GET("/api/forum/:slug/users", handler.GetUsers())
	e.POST("/api/forum/:forumslug/create", handler.CreateThread())
	e.GET("/api/forum/:slug/threads", handler.GetThreads())

	return handler
}

func (delivery *ForumDelivery) CreateForum() echo.HandlerFunc {
	type createForumRequest struct {
		Slug  string `json:"slug" binding:"required"`
		Title string `json:"title" binding:"required"`
		User  string `json:"user" binding:"required"`
	}

	return func(contex echo.Context) error {
		req := &createForumRequest{}
		if err := contex.Bind(req); err != nil {
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		reqForum := &models.Forum{
			Slug:          req.Slug,
			Title:         req.Title,
			AdminNickname: req.User,
		}

		returnForum, err :=delivery.forumUseCase.Add(reqForum)
		if err != nil {
			if err == tools.ErrExistWithSlug {
				return contex.JSON(http.StatusConflict, returnForum)
			}
			if err == tools.ErrUserDoesntExists {
				return contex.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return contex.JSON(http.StatusCreated, returnForum)
	}
}

func (delivery *ForumDelivery) GetForumDetails() echo.HandlerFunc {
	return func(contex echo.Context) error {
		slug := contex.Param("slug")

		returnForum, err := delivery.forumUseCase.GetBySlug(slug)
		if err != nil {
			if err == tools.ErrForumDoesntExists {
				return contex.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return contex.JSON(http.StatusOK, returnForum)
	}
}

func (delivery *ForumDelivery) GetUsers() echo.HandlerFunc {
	return func(contex echo.Context) error {
		slug := contex.Param("slug")
		limit := uint64(0)
		var err error
		if l := contex.QueryParam("limit"); l != "" {
			limit, err = strconv.ParseUint(l, 10, 64)
			if err != nil {
				return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
		}
		since := contex.QueryParam("since")
		desc := false
		if descVal := contex.QueryParam("desc"); descVal == "true" {
			desc = true
		}
		returnUsers, err := delivery.forumUseCase.GetUsers(slug, limit, since, desc)
		if err != nil {
			if err == tools.ErrForumDoesntExists {
				return contex.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return contex.JSON(http.StatusOK, returnUsers)
	}
}

func (delivery *ForumDelivery) CreateThread() echo.HandlerFunc {
		type CreateThreadRequest struct {
		Author  string    `json:"author" binding:"require"`
		Created time.Time `json:"created" binding:"omitempty"`
		Message string    `json:"message" binding:"require"`
		Title   string    `json:"title" binding:"require"`
		Slug    string    `json:"slug" binding:"omitempty"`
	}
	return func(contex echo.Context) error {
		req := &CreateThreadRequest{}
		if err := contex.Bind(req); err != nil {
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		if _, err := strconv.ParseInt(req.Slug, 10, 64); err == nil {
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: tools.ErrIncorrectSlug.Error(),
			})
		}
		slug := contex.Param("forumslug")
		if req.Created.IsZero() {
			req.Created = time.Now()
		}
		reqThread := &models.Thread{
			Author:       req.Author,
			CreationDate: req.Created,
			About:        req.Message,
			Title:        req.Title,
			Slug:         req.Slug,
			Forum:        slug,
		}
		newThread, err := delivery.threadUseCase.CreateThread(reqThread)
		if err != nil {
			if err == tools.ErrForumDoesntExists {
				return contex.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			if err == tools.ErrUserDoesntExists {
				return contex.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			if err == tools.ErrExistWithSlug {
				return contex.JSON(http.StatusConflict, newThread)
			}

			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return contex.JSON(http.StatusCreated, newThread)
	}
}

func (delivery *ForumDelivery) GetThreads() echo.HandlerFunc {
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
		returnThreads, err := delivery.forumUseCase.GetThreadsByForum(slug, limit, since, desc)
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
		return c.JSON(http.StatusOK, returnThreads)
	}
}