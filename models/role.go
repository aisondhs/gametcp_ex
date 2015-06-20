package models

import (
	//"github.com/aisondhs/gametcp_ex/lib/redis"
	//"strconv"
	"time"
)

var Role RoleModel

//
type RoleModel struct {
}

func (this RoleModel) GetRoleByArea(uid int64, srvid int32) (RoleTable, error) {
	var Role RoleTable
	_, err := DataBase().Where("Uid = ? AND Srvid = ? ", uid, srvid).Get(&Role)
	return Role, err
}

func (this RoleModel) Role(roleId int64) (RoleTable, error) {
	var Role RoleTable
	_, err := DataBase().Id(roleId).Get(&Role)
	return Role, err
}

func (this RoleModel) Insert(uid int64, srvid int32) (RoleTable, error) {
	var Role RoleTable
	Role.Uid = uid
	Role.Srvid = srvid
	now := time.Now()
	Role.Ctime = now.Unix()
	_, err := DataBase().Insert(&Role)
	return Role, err
}
