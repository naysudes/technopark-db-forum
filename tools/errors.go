package tools

import "errors"

var (
	ErrDoesntExists           = errors.New("Record doesn't exist")
	ErrUserExistWith          = errors.New("User already exists")
	ErrExistWithSlug          = errors.New("Record exists")
	ErrUserDoesntExists       = errors.New("User doesn't exist")
	ErrForumDoesntExists      = errors.New("Forum with such slug doesn't exist")
	ErrPostDoesntExists       = errors.New("Post doesn't exist")
	ErrParentPostDoesntExists = errors.New("Post doesn't exist")
	ErrThreadDoesntExists     = errors.New("Thread with such slug doesn't exist")
	ErrIncorrectSlug          = errors.New("Slug is incorrect")
	ErrPostIncorrectThreadID  = errors.New("Post has incorrect thread")
)
