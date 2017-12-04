package falcon_portal

type Plugin struct {
	ID         int64  `json:"id" gorm:"column:id"`
	GrpId      int64  `json:"grp_id" gorm:"column:grp_id"`
	Dir        string `json:"dir" gorm:"column:dir"`
	CreateUser string `json:"create_user" gorm:"column:create_user"`
}

func (this Plugin) TableName() string {
	return "plugin_dir"
}

