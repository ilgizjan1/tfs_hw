package users

import "fmt"

type User struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickname"`
	Email    string `json:"email"`
}

func (u *User) String() string {
	return fmt.Sprintf("<User(id=%d, nickName=`%s`, email=`%s`)>", u.ID, u.NickName, u.Email)
}
