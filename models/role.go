package models

import (
	"github.com/aisondhs/gametcp_ex/lib/redis"
	"strconv"
	"time"
)

var Role RoleModel

//
type RoleModel struct {

}

func (this RoleModel) GetRoleByArea(uid int64, areaId int32) (RoleTable, error) {
	var Role RoleTable
	_, err := DataBase().Where("uid = ? AND area_id = ? ", uid, areaId).Get(&Role)
	return Role, err
}

func (this RoleModel) Role(roleId int64) (RoleTable, error) {
	var Role RoleTable
	_, err := DataBase().Id(roleId).Get(&Role)
	return Role, err
}

func (this RoleModel) Insert(uid int64, areaId int32) (RoleTable, error) {
	var Role RoleTable
	Role.Uid = uid
	Role.Areaid = areaId
	now := time.Now()
	Role.Ctime = now.Unix()
	_, err := DataBase().Insert(&Role)
	return Role, err
}

// 设置最新登录角色
func (this RoleModel) SetLastRole(Role RoleTable) error {
	return redis.Redis.Set("UID_RID_"+strconv.FormatInt(Role.Uid, 10), strconv.FormatInt(Role.RoleId, 10))
}

// 取得最新登录角色
func (this RoleModel) GetLastRole(uid int64) (int64, error) {
	if s, err := redis.Redis.Get("UID_RID_" + strconv.FormatInt(uid, 10)); err != nil {
		return 0, err
	} else {
		return strconv.ParseInt(s, 10, 0)
	}
}
