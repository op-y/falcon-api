package dashboard

type DashboardScreen struct {
	ID   int64  `json:"id" gorm:"column:id"`
	PID  int64  `json:"pid" gorm:"column:pid"`
	Name string `json:"name" gorm:"column:name"`
}

func (this DashboardScreen) TableName() string {
	return "dashboard_screen"
}

