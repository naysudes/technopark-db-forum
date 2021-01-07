package tools

import "errors"

var (
	ErrDoesntExists           = errors.New("Record doesn't exists")
	ErrUserExistWith          = errors.New("User already exists")
	ErrExistWithSlug          = errors.New("Record exists")
	ErrUserDoesntExists       = errors.New("User doesn't exists")
	ErrForumDoesntExists      = errors.New("Forum doesn't exists")
	ErrPostDoesntExists       = errors.New("Post doesn't exists")
	ErrParentPostDoesntExists = errors.New("Post doesn't exists")
	ErrThreadDoesntExists     = errors.New("Thread doesn't exists")
	ErrIncorrectSlug          = errors.New("Slug is incorrect")
	ErrPostIncorrectThreadID  = errors.New("Post has incorrect thread")
)
