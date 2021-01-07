package forum

import "github.com/naysudes/technopark-db-forum/models"

type Usecase interface {
	AddForum(*models.Forum) (*models.Forum, error)
	GetForumBySlug(string) (*models.Forum, error)
	GetForumUsers(string, uint64, string, bool) ([]*models.User, error)
	GetForumThreads(string, uint64, string, bool) ([]*models.Thread, error)
}
