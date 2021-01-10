package models

type Status struct {
	UsersCount   uint64 `json:"user"`
	ThreadsCount uint64 `json:"thread"`
	ForumsCount  uint64 `json:"forum"`
	PostsCount   uint64 `json:"post"`
}