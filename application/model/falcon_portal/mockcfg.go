package falcon_portal

//no_data
type Mockcfg struct {
	ID   int64  `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
	Obj  string `json:"obj" gorm:"column:obj"`
	//group, host, other
	ObjType string  `json:"obj_type" gorm:"column:obj_type"`
	Metric  string  `json:"metric" gorm:"column:metric"`
	Tags    string  `json:"tags" gorm:"column:tags"`
	DsType  string  `json:"dstype" gorm:"column:dstype"`
	Step    int     `json:"step" gorm:"column:step"`
	Mock    float64 `json:"mock" gorm:"column:mock"`
	Creator string  `json:"creator" gorm:"column:creator"`
}

func (this Mockcfg) TableName() string {
	return "mockcfg"
}

