package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/naysudes/technopark-db-forum/models"
	// "github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/post"
)

type PostRepository struct {
	db *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) post.Repository {
	return  PostRepository{
		db: db,
	}
}

func (repo PostRepository) Insert(posts []*models.Post) error {
		queryRow := "INSERT INTO posts (author, forum, message, parent, thread) VALUES "
	for _, p := range posts {
		queryRow += fmt.Sprintf("(%d, %d, '%s', %d, %d),",
			p.AuthorID, p.ForumID, p.Message, p.ParentID, p.ThreadID)
		_, err := repo.db.Exec("INSERT INTO forums_users (user_id, forum_id) VALUES "+
			"($1, $2) ON CONFLICT DO NOTHING", p.AuthorID, p.ForumID)
		if err != nil {
			return err
		}
	}
	queryRow = queryRow[0 : len(queryRow)-1]
	queryRow += " RETURNING id, created"

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	rows, err := tx.Query(queryRow)
	if err != nil {
		return err
	}
	defer rows.Close()
	postIndex := 0
	for rows.Next() {
		if err := rows.Scan(&posts[postIndex].ID, &posts[postIndex].CreationDate); err != nil {
			return err
		}
		postIndex++
	}
	return tx.Commit()
}

func (repo PostRepository) GetByThread(id uint64, limit uint64, since uint64, sort string, desc bool) ([]*models.Post, error) {
	queryStr := "SELECT p.id, u.nickname, f.slug, p.thread, p.created, p.message, p.isEdited, " +
		"coalesce(p.path[array_length(p.path, 1) - 1], 0) " +
		"FROM posts AS p " +
		"JOIN users AS u ON (u.id = p.author) " +
		"JOIN forums AS f ON (f.id = p.forum) " +
		"WHERE %s %s"

	var where string
	var order string
	switch sort {
	case "flat", "":
		where = "p.thread = $1"
		if since != 0 {
			if desc {
				where += " AND p.id < $2"
			} else {
				where += " AND p.id > $2"
			}
		}
		order = "ORDER BY "
		if sort == "flat" {
			order += "p.created"
			if desc {
				order += " DESC"
			}
			order += ", p.id"
			if desc {
				order += " DESC"
			}
		} else {
			order += "p.id"
			if desc {
				order += " DESC"
			}
		}
		if limit != 0 {
			if since != 0 {
				order += " LIMIT $3"
			} else {
				order += " LIMIT $2"
			}
		}
	case "tree":
		where = "p.thread = $1"
		if since != 0 {
			if desc {
				where += " AND coalesce(path < (select path FROM posts where id = $2), true)"
			} else {
				where += " AND coalesce(path > (select path FROM posts where id = $2), true)"
			}
		}
		order = "ORDER BY p.path[1]"
		if desc {
			order += " DESC"
		}
		order += ", p.path[2:]"
		if desc {
			order += " DESC"
		}
		order += " NULLS FIRST"
		if limit != 0 {
			if since != 0 {
				order += " LIMIT $3"
			} else {
				order += " LIMIT $2"
			}
		}
	case "parent_tree":
		where = "p.path[1] IN (SELECT path[1] FROM posts WHERE thread = $1 AND " +
			"array_length(path, 1) = 1"
		if since != 0 {
			if desc {
				where += " AND id < (SELECT path[1] FROM posts WHERE id = $2)"
			} else {
				where += " AND id > (SELECT path[1] FROM posts WHERE id = $2)"
			}
		}
		where += " ORDER BY id"
		if desc {
			where += " DESC"
		}
		if limit != 0 {
			if since != 0 {
				where += " LIMIT $3"
			} else {
				where += " LIMIT $2"
			}
		}
		where += ")"
		order = "ORDER BY p.path[1]"
		if desc {
			order += " DESC"
		}
		order += ", p.path[2:] NULLS FIRST"
	}

	rows := &pgx.Rows{}
	var err error
	if since != 0 {
		if limit != 0 {
			rows, err = repo.db.Query(fmt.Sprintf(queryStr, where, order), id, since, limit)
		} else {
			rows, err = repo.db.Query(fmt.Sprintf(queryStr, where, order), id, since)
		}
	} else {
		if limit != 0 {
			rows, err = repo.db.Query(fmt.Sprintf(queryStr, where, order), id, limit)
		} else {
			rows, err = repo.db.Query(fmt.Sprintf(queryStr, where, order), id)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	postsByThread := []*models.Post{}
	for rows.Next() {
		post := &models.Post{}

		err = rows.Scan(&post.ID, &post.Author, &post.Forum, &post.ThreadID,
			&post.CreationDate, &post.Message, &post.IsEdited, &post.ParentID)
		if err != nil {
			return nil, err
		}
		postsByThread = append(postsByThread, post)
	}
	return postsByThread, nil
}
func (repo PostRepository) CheckParentPosts([]*models.Post, uint64) (bool, error) {
	return true, nil
}
func (repo PostRepository) GetByID(uint64) (*models.Post, error) {
	return nil, nil
}
func (repo PostRepository) GetCountByForumID(uint64) (uint64, error) {
	return 1, nil
}
func (repo PostRepository) Update(*models.Post) error {
	return nil
}
