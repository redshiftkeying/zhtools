package user

import (
	"context"
	"github.com/redshiftkeying/zhtools/user/repo"
)

type User struct {
	ctx  context.Context
	repo repo.Queries
}

type Profile struct {
	Status int
}

func (u *User) Login(username, password string) (bool, error) {
	res, err := u.repo.GetPasswordByName(u.ctx, username)
	if err != nil || !res.Valid {
		return false, err
	}
	return res.String == password, err
}
