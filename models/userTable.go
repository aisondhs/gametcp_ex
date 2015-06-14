package models

// user
type UserTable struct {
	Uid          int64  `xorm:"NOT NULL pk autoincr INT(11)"`
	Account   string `xorm:"NOT NULL DEFAULT '' VARCHAR(64)"`
	Pwd string `xorm:"NOT NULL DEFAULT '' VARCHAR(64)"`
	Ctime int64 `xorm:"NOT NULL DEFAULT 0 INT(11)"`
}

func (this UserTable) TableName() string {
	return "user"
}