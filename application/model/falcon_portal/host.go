package falcon_portal

import (
	"errors"
	"fmt"

	"falcon-api/application/model"
)

type Host struct {
	ID            int64  `json:"id" gorm:"column:id"`
	Hostname      string `json:"hostname" gorm:"column:hostname"`
	Ip            string `json:"ip" gorm:"column:ip"`
	AgentVersion  string `json:"agent_version"  gorm:"column:agent_version"`
	PluginVersion string `json:"plugin_version"  gorm:"column:plugin_version"`
	MaintainBegin int64  `json:"maintain_begin"  gorm:"column:maintain_begin"`
	MaintainEnd   int64  `json:"maintain_end"  gorm:"column:maintain_end"`
}

func (this Host) TableName() string {
	return "host"
}

func (this Host) Existing() (int64, bool) {
	db := model.Con()
	db.FalconPortal.Table(this.TableName()).Where("hostname = ?", this.Hostname).Scan(&this)
	if this.ID != 0 {
		return this.ID, true
	} else {
		return 0, false
	}
}

func (this Host) RelatedGrp() (Grps []HostGroup) {
	db := model.Con()
	grpHost := []GrpHost{}
	db.FalconPortal.Select("grp_id").Where("host_id = ?", this.ID).Find(&grpHost)
	tids := []int64{}
	for _, t := range grpHost {
		tids = append(tids, t.GrpID)
	}
	tidStr, _ := arrInt64ToString(tids)
	Grps = []HostGroup{}
	db.FalconPortal.Where(fmt.Sprintf("id in (%s)", tidStr)).Find(&Grps)
	return
}

func (this Host) RelatedTpl() (tpls []Template) {
	db := model.Con()
	grps := this.RelatedGrp()
	gids := []int64{}
	for _, g := range grps {
		gids = append(gids, g.ID)
	}
	gidStr, _ := arrInt64ToString(gids)
	grpTpls := []GrpTpl{}
	db.FalconPortal.Select("tpl_id").Where(fmt.Sprintf("grp_id in (%s)", gidStr)).Find(&grpTpls)
	tids := []int64{}
	for _, t := range grpTpls {
		tids = append(tids, t.TplID)
	}
	tidStr, _ := arrInt64ToString(tids)
	tpls = []Template{}
	db.FalconPortal.Where(fmt.Sprintf("id in (%s)", tidStr)).Find(&tpls)
	return
}

func arrInt64ToString(arr []int64) (result string, err error) {
	result = ""
	for indx, a := range arr {
		if indx == 0 {
			result = fmt.Sprintf("%v", a)
		} else {
			result = fmt.Sprintf("%v,%v", result, a)
		}
	}
	if result == "" {
		err = errors.New(fmt.Sprintf("array is empty, err: %v", arr))
	}
	return
}

