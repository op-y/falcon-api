package falcon_portal

import (
	"falcon-api/application/model"
)

type Cluster struct {
	ID          int64  `json:"id" gorm:"column:id"`
	GrpId       int64  `json:"grp_id" gorm:"column:grp_id"`
	Numerator   string `json:"numerator" gorm:"column:numerator"`
	Denominator string `json:"denominator" gorm:"denominator"`
	Endpoint    string `json:"endpoint" gorm:"endpoint"`
	Metric      string `json:"metric" gorm:"metric"`
	Tags        string `json:"tags" gorm:"tags"`
	DsType      string `json:"ds_type" gorm:"ds_type"`
	Step        int    `json:"step" gorm:"step"`
	Creator     string `json:"creator" gorm:"creator"`
}

func (this Cluster) TableName() string {
	return "cluster"
}

func (this Cluster) HostGroupName() (name string, err error) {
	if this.GrpId == 0 {
		return
	}
	db := model.Con()
	var hg HostGroup
	hg.ID = this.GrpId
	if dt := db.FalconPortal.Find(&hg); dt.Error != nil {
		return name, dt.Error
	}
	name = hg.Name
	return
}

