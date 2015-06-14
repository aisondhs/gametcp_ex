package models

import ()

func init() {

}

var User UserModel

type UserModel struct {
}

func (this UserModel) Insert(u *UserTable) error {
	_, err := DataBase().Insert(u)
	return err
}

func (this UserModel) User(uid int64) (*UserTable, error) {
	u := new(UserTable)
	return u, DataBase().Id(uid).Find(u)
}

func (this UserModel) GetUserByAccount(account string) (*UserTable, error) {
	u := new(UserTable)
	_, err := DataBase().Where("Account = ? ", account).Get(u)
	return u, err
}