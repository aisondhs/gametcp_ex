package models

type RoleTable struct {
	RoleId int64 `xorm:"pk autoincr 'RoleId'"`
	Uid    int64
	Srvid  int32
	Ctime  int64
}

func (this RoleTable) TableName() string {
	return "role"
}

func (this *RoleTable) Update() error {
	_, err := DataBase().Id(this.RoleId).Update(this)
	return err
}
