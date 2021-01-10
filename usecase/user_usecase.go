package usecase

import (
	"github.com/naysudes/technopark-db-forum/models"
	"github.com/naysudes/technopark-db-forum/tools"
	"github.com/naysudes/technopark-db-forum/interfaces/user"
)

type UserUsecase struct {
	repo user.Repository
}

func NewUserUsecase(repo user.Repository) user.Usecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (usecase *UserUsecase) Add(nickname string, usr *models.User) ([]*models.User, error) {
	usr1, err := usecase.repo.GetByNickname(nickname)
	if err != nil && err != tools.ErrDoesntExists {
		return nil, err
	}
	usr2, err := usecase.repo.GetByEmail(usr.Email)
	if err != nil && err != tools.ErrDoesntExists {
		return nil, err
	}

	if usr1 != nil || usr2 != nil {
		returnUsers := []*models.User{}
		if usr1 != nil {
			returnUsers = append(returnUsers, usr1)
			if usr2 != nil && usr1.Nickname != usr2.Nickname {
				returnUsers = append(returnUsers, usr2)
			}
		} else if usr2 != nil {
			returnUsers = append(returnUsers, usr2)
		}
		return returnUsers, tools.ErrUserExistWith
	}

	usr.SetNickname(nickname)
	if err = usecase.repo.Insert(usr); err != nil {
		return nil, err
	}

	return []*models.User{usr}, nil
}

func (usecase *UserUsecase) GetByNickname(nickname string) (*models.User, error) {
	user, err := usecase.repo.GetByNickname(nickname)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (usecase *UserUsecase) Update(nickname string, usr *models.User) error {
	user, err := usecase.repo.GetByNickname(nickname)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return tools.ErrUserDoesntExists
		}
		return err
	}
	usrCheckEmail, err := usecase.repo.GetByEmail(usr.Email)
	if err != nil && err != tools.ErrDoesntExists {
		return err
	}
	if err != tools.ErrDoesntExists && usrCheckEmail.Nickname != user.Nickname {
		return tools.ErrUserExistWith
	}
	usr.SetNickname(nickname)
	if usr.About == "" {
		usr.About = user.About
	}
	if usr.Email == "" {
		usr.Email = user.Email
	}
	if usr.Fullname == "" {
		usr.Fullname = user.Fullname
	}
	if err = usecase.repo.Update(usr); err != nil {
		return err
	}
	return nil
}
