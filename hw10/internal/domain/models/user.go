package models

import "fmt"

type User struct {
	ID       int64
	NickName string
	Email    string
}

func (u *User) String() string {
	return fmt.Sprintf("<User(id=%d, nickName=`%s`, email=`%s`)>", u.ID, u.NickName, u.Email)
}
