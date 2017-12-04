package uic

import (
    "falcon-api/application/model"
)

type RelTeamUser struct {
	ID  int64
	Tid int64
	Uid int64
}

func (this RelTeamUser) TableName() string {
	return "rel_team_user"
}

func (this RelTeamUser) Me() {
	db := model.Con()
	db.Uic.Where("id = 1")
}

