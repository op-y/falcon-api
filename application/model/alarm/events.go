package alarm

import "time"

type Events struct {
	ID          int64      `json:"id" gorm:"column:id"`
	EventCaseId string     `json:"event_caseId" gorm:"column:event_caseId"`
	Step        int        `json:"step" gorm:"step"`
	Cond        string     `json:"cond" gorm:"cond"`
	Status      int        `json:"status" gorm:"status"`
	Timestamp   *time.Time `json:"timestamp" gorm:"timestamp"`
}

func (this Events) TableName() string {
	return "events"
}
