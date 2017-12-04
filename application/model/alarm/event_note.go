package alarm

import (
	"time"

    "falcon-api/application/model"
    "falcon-api/application/model/uic"
)

type EventNote struct {
	ID          int64      `json:"id" gorm:"column:id"`
	EventCaseId string     `json:"event_caseId" gorm:"column:event_caseId"`
	Note        string     `json:"note" gorm:"note"`
	CaseId      string     `json:"case_id" gorm:"case_id"`
	Status      string     `json:"status" gorm:"status"`
	Timestamp   *time.Time `json:"timestamp" gorm:"timestamp"`
	UserId      int64      `json:"user_id" gorm:"user_id"`
}

func (this EventNote) TableName() string {
	return "event_note"
}

func (this EventNote) GetUserName() string {
	db := model.Con()
	user := uic.User{ID: this.UserId}
	db.Uic.Table(user.TableName()).Where(&user).Scan(&user)
	return user.Name
}
