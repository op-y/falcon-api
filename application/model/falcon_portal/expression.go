package falcon_portal

type Expression struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Expression string `json:"expression" gorm:"column:expression"`
	Func       string `json:"func" gorm:"column:func"`
	Op         string `json:"op" gorm:"column:op"`
	RightValue string `json:"right_value" gorm:"column:right_value"`
	MaxStep    int    `json:"max_step" gorm:"column:max_step"`
	Priority   int    `json:"priority" gorm:"column:priority"`
	Note       string `json:"note" gorm:"column:note"`
	ActionId   int64  `json:"action_id" gorm:"column:action_id"`
	CreateUser string `json:"create_user" gorm:"column:create_user"`
	Pause      int    `json:"pause" gorm:"column:pause"`
}

