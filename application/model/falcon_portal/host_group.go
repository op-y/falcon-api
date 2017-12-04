package falcon_portal

type HostGroup struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Name       string `json:"grp_name" gorm:"column:grp_name"`
	CreateUser string `json:"create_user" gorm:"column:create_user"`
	ComeFrom   int    `json:"-"  gorm:"column:come_from"`
}

func (this HostGroup) TableName() string {
	return "grp"
}

