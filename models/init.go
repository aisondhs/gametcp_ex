package models

import (
	"github.com/aisondhs/gametcp_ex/lib/db"
	"github.com/go-xorm/xorm"
)

var Db *xorm.Engine

func DataBase() *xorm.Engine {
	if Db == nil {
		Db = db.DataBase
	}
	return Db
}
