package dashboard

type DashboardTmpGraph struct {
	ID        int64  `json:"id" gorm:"column:id"`
	Endpoints string `json:"endpoints" gorm:"column:endpoints"`
	Counters  string `json:"counters" gorm:"column:counters"`
	CK        string `json:"ck" gorm:"column:ck"`
}

func (this DashboardTmpGraph) TableName() string {
	return "tmp_graph"
}

