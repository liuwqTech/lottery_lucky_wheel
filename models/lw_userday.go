package models

type LwUserday struct {
	Id         int `xorm:"not null pk autoincr comment('主键') INT(11)"`
	Uid        int `xorm:"comment('用户id') INT(11)"`
	Day        int `xorm:"comment('日期，如20180725') INT(11)"`
	Num        int `xorm:"comment('次数') INT(11)"`
	SysCreated int `xorm:"comment('创建时间') INT(11)"`
	SysUpdated int `xorm:"comment('修改时间') INT(11)"`
}
