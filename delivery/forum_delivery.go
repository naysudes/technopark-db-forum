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

type ForumHandler struct {
	forumUseCase  forum.Usecase
	threadUseCase thread.Usecase
}

func NewForumHandler(e *echo.Echo, thUC thread.Usecase, fUC forum.Usecase) ForumHandler {
	handler := ForumHandler{ forumUseCase:  fUC, threadUseCase: thUC}

	e.POST("/api/forum/create", handler.CreateForum())
	e.GET("/api/forum/:slug/details", handler.GetForumDetails())
	e.GET("/api/forum/:slug/users", handler.GetUsers())
	e.POST("/api/forum/:forumslug/create", handler.CreateThread())

	return handler
}

func (handler ForumHandler) CreateForum() echo.HandlerFunc {
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

		returnForum, err := handler.forumUseCase.Add(reqForum)
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

func (handler ForumHandler) GetForumDetails() echo.HandlerFunc {
	return func(contex echo.Context) error {
		slug := contex.Param("slug")

		returnForum, err := handler.forumUseCase.GetBySlug(slug)
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

func (handler ForumHandler) GetUsers() echo.HandlerFunc {
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

		returnUsers, err := handler.forumUseCase.GetUsers(slug, limit, since, desc)
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

func (handler ForumHandler) CreateThread() echo.HandlerFunc {
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

		newThread, err := handler.threadUseCase.CreateThread(reqThread)
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