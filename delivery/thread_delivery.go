package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/forum"
	"github.com/naysudes/technopark-db-forum/interfaces/thread"
	"net/http"
	"strconv"
	"encoding/json"
)

type ThreadDelivery struct {
	forumUseCase  forum.Usecase
	threadUseCase thread.Usecase
}

func NewThreadDelivery(e *echo.Echo, forumUC forum.Usecase, threadUC thread.Usecase) *ThreadDelivery {
	handler := &ThreadDelivery{ forumUseCase:  forumUC, threadUseCase: threadUC }
	e.POST("/api/thread/:slug_or_id/create", handler.CreatePosts())
	e.GET("/api/thread/:slug_or_id/details", handler.GetDetails())
	e.POST("/api/thread/:slug_or_id/details", handler.Update())
	e.GET("/api/thread/:slug_or_id/posts", handler.GetPosts())
	e.POST("/api/thread/:slug_or_id/vote", handler.Vote())
	return handler
}

func (delivery *ThreadDelivery) CreatePosts() echo.HandlerFunc {
	type createReq struct {
		Author  string `json:"author" binding:"required"`
		Message string `json:"message" binding:"required"`
		Parent  uint64 `json:"parent" binding:"required"`
	}
  	return func(contex echo.Context) error {
		req := []createReq{}
		if err := json.NewDecoder(contex.Request().Body).Decode(&req); err != nil {
			return contex.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		slugOrID := contex.Param("slug_or_id")

		posts := make([]*models.Post, 0, len(req))
		for _, r := range req {
			post := &models.Post{
				Author:   r.Author,
				Message:  r.Message,
				ParentID: r.Parent,
				IsEdited: false,
			}
			posts = append(posts, post)
		}
		returnPosts, err := delivery.threadUseCase.CreatePosts(slugOrID, posts)
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

func (delivery *ThreadDelivery) GetDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		slugOrID := c.Param("slug_or_id")
		threadBySlugorId, err := delivery.threadUseCase.GetBySlugOrID(slugOrID)
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

func (delivery *ThreadDelivery) Update() echo.HandlerFunc {
	type updateThreadRequest struct {
		Message string `json:"message" binding:"require"`
		Title   string `json:"title" binding:"require"`
	}
	return func(c echo.Context) error {
		req := &updateThreadRequest{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}

		slugOrID := c.Param("slug_or_id")
		reqThread := &models.Thread{
			About: req.Message,
			Title: req.Title,
		}

		returnThread, err := delivery.threadUseCase.Update(slugOrID, reqThread)
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

		return c.JSON(http.StatusOK, returnThread)
	}
}

func (delivery *ThreadDelivery) GetPosts() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrId := context.Param("slug_or_id")
		limit := uint64(0)
		var err error
		if l := context.QueryParam("limit"); l != "" {
			limit, err = strconv.ParseUint(l, 10, 64)
			if err != nil {
				return context.JSON(http.StatusBadRequest, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
		}
		since := uint64(0)
		if s := context.QueryParam("since"); s != "" {
			since, err = strconv.ParseUint(s, 10, 64)
			if err != nil {
				return context.JSON(http.StatusBadRequest, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
		}
		desc := false
		if descVal := context.QueryParam("desc"); descVal == "true" {
			desc = true
		}
		sort := context.QueryParam("sort")
		threadsByForum, err := delivery.threadUseCase.GetPosts(slugOrId, limit, since, sort, desc)
		if err != nil {
			if err == tools.ErrThreadDoesntExists {
				return context.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}
			return context.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return context.JSON(http.StatusOK, threadsByForum)
	}
}

func (delivery *ThreadDelivery) Vote() echo.HandlerFunc{
		type voteReq struct {
		Nickname string `json:"nickname" binding:"require"`
		Voice    int64  `json:"voice" binding:"require"`
	}
	return func(c echo.Context) error {
		req := &voteReq{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		slugOrID := c.Param("slug_or_id")
		vReq := &models.Vote{
			Nickname: req.Nickname,
			Voice:    req.Voice,
		}
		thread, err := delivery.threadUseCase.Vote(slugOrID, vReq)
		if err != nil {
			if err == tools.ErrThreadDoesntExists || err == tools.ErrUserDoesntExists {
				return c.JSON(http.StatusNotFound, tools.ErrorResponce{
					Message: err.Error(),
				})
			}

			return c.JSON(http.StatusBadRequest, tools.ErrorResponce{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, thread)
	}
}
