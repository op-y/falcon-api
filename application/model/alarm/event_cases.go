package alarm

import (
	"fmt"
	"time"

    "falcon-api/application/model"
)

type EventCases struct {
	ID            string     `json:"id" gorm:"column:id"`
	Endpoint      string     `json:"endpoint" gorm:"column:endpoint"`
	Metric        string     `json:"metric" gorm:"metric"`
	Func          string     `json:"func" gorm:"func"`
	Cond          string     `json:"cond" gorm:"cond"`
	Note          string     `json:"note" gorm:"note"`
	MaxStep       int        `json:"step" gorm:"step"`
	CurrentStep   int        `json:"current_step" gorm:"current_step"`
	Priority      int        `json:"priority" gorm:"priority"`
	Status        string     `json:"status" gorm:"status"`
	Timestamp     *time.Time `json:"timestamp" gorm:"timestamp"`
	UpdateAt      *time.Time `json:"update_at" gorm:"update_at"`
	ClosedAt      *time.Time `json:"closed_at" gorm:"closed_at"`
	ClosedNote    string     `json:"closed_note" gorm:"closed_note"`
	UserModified  int64      `json:"user_modified" gorm:"user_modified"`
	TplCreator    string     `json:"tpl_creator" gorm:"tpl_creator"`
	ExpressionId  int64      `json:"expression_id" gorm:"expression_id"`
	StrategyId    int64      `json:"strategy_id" gorm:"strategy_id"`
	TemplateId    int64      `json:"template_id" gorm:"template_id"`
	ProcessNote   int64      `json:"process_note" gorm:"process_note"`
	ProcessStatus string     `json:"process_status" gorm:"process_status"`
}

func (this EventCases) TableName() string {
	return "event_cases"
}

func (this EventCases) GetEvents() []Events {
	db := model.Con()
	t := Events{
		EventCaseId: this.ID,
	}
	e := []Events{}
	db.Alarm.Table(t.TableName()).Where(&t).Scan(&e)
	return e
}

func (this EventCases) GetNotes() []EventNote {
	db := model.Con()
	perpareSql := fmt.Sprintf("event_caseId = '%s' AND timestamp >= FROM_UNIXTIME(%d)", this.ID, this.Timestamp.Unix())
	t := EventCases{}
	notes := []EventNote{}
	db.Alarm.Table(t.TableName()).Where(perpareSql).Scan(&notes)
	return notes
}

func (this EventCases) NotesCount() int {
	notes := this.GetNotes()
	return len(notes)
}
