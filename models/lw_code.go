package models

type LwCode struct {
	Id         int    `xorm:"not null pk autoincr comment('主键') INT(11)"`
	GiftId     int    `xorm:"comment('奖品id，关联lw_gift表') INT(11)"`
	Code       string `xorm:"comment('虚拟券编码') VARCHAR(255)"`
	SysCreated int    `xorm:"comment('创建时间') INT(11)"`
	SysUpdated int    `xorm:"comment('更新时间') INT(11)"`
	SysStatus  int    `xorm:"comment('状态，0正常，1作废，2已发放') SMALLINT(6)"`
}
