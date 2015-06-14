package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/Unknwon/goconfig"
	//_ "github.com/lib/pq"
)

var DataBase *xorm.Engine

func init() {

	c, err := goconfig.LoadConfigFile("conf/conf.ini")
	if err != nil {
		panic(err)
	}

	driver, err := c.GetValue("Database", "driver")
	if err != nil {
		panic(err)
	}
	dsn, err := c.GetValue("Database", "dsn")

	if err != nil {
		panic(err)
	}

	DataBase, err = xorm.NewEngine(driver, dsn)
	if err != nil {
		panic(err)
	}

	err = DataBase.Ping()
	if err != nil {
		panic(err)
	}
}