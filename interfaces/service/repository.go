package service

type Repository interface {
	GetUserCount() (uint64, error)
	GetThreadCount() (uint64, error)
	GetForumCount() (uint64, error)
	GetPostCount() (uint64, error)

	DeleteUsers() error
	DeleteForums() error
	DeleteThreads() error
	DeletePosts() error
	DeleteVotes() error
}