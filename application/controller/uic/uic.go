package uic

import (
    "fmt"
    "net/http"
    "strconv"
    "time"

    h "falcon-api/application/helper"
    "falcon-api/application/model"
    um "falcon-api/application/model/uic"
    "falcon-api/application/utils"

    log "github.com/Sirupsen/logrus"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
)

type APIUserInput struct {
    Name   string `json:"name" binding:"required"`
    Cnname string `json:"cnname" binding:"required"`
    Passwd string `json:"password" binding:"required"`
    Email  string `json:"email" binding:"required"`
    Phone  string `json:"phone"`
    IM     string `json:"im"`
    QQ     string `json:"qq"`
}

type APIChangePassword struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type APIUserDeleteInput struct {
    Username string `json:"username" binding:"required"`
}

type CTeam struct {
    Team        um.Team   `json:"team"`
    TeamCreator string    `json:"creator_name"`
    Users       []um.User `json:"users"`
}

type APIGetTeamOutput struct {
    um.Team
    Users       []um.User `json:"users"`
    TeamCreator string    `json:"creator_name"`
}

type APICreateTeamInput struct {
    Name    string  `json:"team_name" binding:"required"`                                                
    Resume  string  `json:"resume"`
    UserIDs []int64 `json:"users"`
}

type APIUpdateTeamInput struct {
    ID      int    `json:"team_id" binding:"required"`
    Resume  string `json:"resume"`
    Name    string `json:"name"`
    UserIDs []int  `json:"users"`
}

type APIDeleteTeamInput struct {
    ID int64 `json:"team_id" binding:"required"`
}

var db model.DBPool

func GetUsers(c *gin.Context) {
    var (
        limit int 
        page  int 
        err   error
    )   
    pageTmp := c.DefaultQuery("page", "") 
    limitTmp := c.DefaultQuery("limit", "") 
    page, limit, err = h.PageParser(pageTmp, limitTmp)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err.Error())
        return
    }   
    q := c.DefaultQuery("q", ".+")
    var users []um.User
    var dt *gorm.DB
    if limit != -1 && page != -1 {
        dt = db.Uic.Raw(
            fmt.Sprintf("select * from user where name regexp '%s' limit %d,%d", q, page, limit)).Scan(&users)
    } else {
        dt = db.Uic.Table("user").Where("name regexp ?", q).Scan(&users)
    }   
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }   
    h.JSONR(c, users)
}

func GetUserByName(c *gin.Context) {
    name := c.Params.ByName("name")                                                                 
    if name == "" {
        h.JSONR(c, http.StatusBadRequest, "user name is missing")
        return
    }
    user := um.User{}
    if dt := db.Uic.Table("user").Where("name = ?", name).First(&user); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }
    h.JSONR(c, user)
    return
}

func CreateUser(c *gin.Context) {
    var inputs APIUserInput
    err := c.Bind(&inputs)

    switch {
    case err != nil:
        h.JSONR(c, http.StatusBadRequest, err)
        return
    case utils.HasDangerousCharacters(inputs.Cnname):
        h.JSONR(c, http.StatusBadRequest, "name pattern is invalid")
        return
    }

    var user um.User
    db.Uic.Table("user").Where("name = ?", inputs.Name).Scan(&user)
    if user.ID != 0 {
        h.JSONR(c, http.StatusBadRequest, "name is already existing")
        return
    }
    password := utils.HashIt(inputs.Passwd)
    user = um.User{
        Name:   inputs.Name,
        Passwd: password,
        Cnname: inputs.Cnname,
        Email:  inputs.Email,
        Phone:  inputs.Phone,
        IM:     inputs.IM,
        QQ:     inputs.QQ,
    }

    //for create a root user during the first time
    if inputs.Name == "root" {
        user.Role = 2
    }

    dt := db.Uic.Table("user").Create(&user)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    var session um.Session
    response := map[string]string{}
    s := db.Uic.Table("session").Where("uid = ?", user.ID).Scan(&session)
    if s.Error != nil && s.Error.Error() != "record not found" {
        h.JSONR(c, http.StatusBadRequest, s.Error)
        return
    } else if session.ID == 0 {
        session.Sig = utils.GenerateUUID()
        session.Expired = int(time.Now().Unix()) + 3600*24*30
        session.Uid = user.ID
        db.Uic.Create(&session)
    }
    log.Debugf("%v", session)
    response["sig"] = session.Sig
    response["name"] = user.Name
    h.JSONR(c, http.StatusOK, response)
    return
}

func UpdateUserProfile(c *gin.Context) {
    var inputs APIUserInput
    err := c.BindJSON(&inputs)
    if err != nil {
        h.JSONR(c, http.StatusExpectationFailed, err)
        return
    }

    user := um.User{}
    name := inputs.Name
    uuser := map[string]interface{}{
        "Cnname": inputs.Cnname,
        "Email":  inputs.Email,
        "Phone":  inputs.Phone,
        "IM":     inputs.IM,
        "QQ":     inputs.QQ,
    }
    dt := db.Uic.Model(&user).Where("name = ?", name).Update(uuser)
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }
    h.JSONR(c, "user profile updated")
    return
}

func UpdateUserPassword(c *gin.Context) {
    var inputs APIChangePassword
    err := c.Bind(&inputs)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    user := um.User{Name: inputs.Username}
    dt := db.Uic.Where(&user).Find(&user)
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }

    user.Passwd = utils.HashIt(inputs.Password)
    dt = db.Uic.Save(&user)
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }
    h.JSONR(c, http.StatusOK, "password updated!")
    return
}

func GetUserTeams(c *gin.Context) {
    username := c.Params.ByName("name")
    if username == "" {
        h.JSONR(c, http.StatusBadRequest, "user name is missing")
        return
    }

    user := um.User{}
    dt := db.Uic.Table("user").Where("name = ?", username).First(&user)
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, "get user fail")
        return
    }

    tus := []um.RelTeamUser{}
    dt = db.Uic.Table("rel_team_user").Where("uid = ?", user.ID).Find(&tus)
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, "get team ids fail")
        return
    }
    tids := []int64{}
    for _, ut := range tus {
        tids = append(tids, ut.Tid)
    }
    teams := []um.Team{}
    tidsStr, _ := utils.ArrInt64ToString(tids)
    // BUG:: exception was thrown when tids is nil
    dt = db.Uic.Table("team").Where(fmt.Sprintf("id in (%s)", tidsStr)).Find(&teams)
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, "get teams fail")
        return
    }
    h.JSONR(c, map[string]interface{}{
        "teams": teams,
    })
    return
}

func DeleteUser(c *gin.Context) {
    var inputs APIUserDeleteInput
    err := c.Bind(&inputs)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    } 

    dt := db.Uic.Where("name = ?", inputs.Username).Delete(&um.User{})
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    } else if dt.RowsAffected == 0 { 
        h.JSONR(c, http.StatusExpectationFailed, "you have no such permission or sth goes wrong")
        return
    }
    h.JSONR(c, fmt.Sprintf("user %v has been delete, affect row: %v", inputs.Username, dt.RowsAffected))
    return
}

func GetTeams(c *gin.Context) {
    var (
        limit int 
        page  int 
        err   error
    )   
    pageTmp := c.DefaultQuery("page", "") 
    limitTmp := c.DefaultQuery("limit", "") 
    page, limit, err = h.PageParser(pageTmp, limitTmp)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err.Error())
        return
    }   
    query := c.DefaultQuery("q", ".+")

    var dt *gorm.DB
    teams := []um.Team{}
    if limit != -1 && page != -1 {
        dt = db.Uic.Table("team").Raw(
            "select * from team where name regexp ? limit ?, ?", query, page, limit).Scan(&teams)
    } else {
        dt = db.Uic.Table("team").Where("name regexp ?", query).Scan(&teams)
    }   
    err = dt.Error
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    outputs := []CTeam{}
    for _, t := range teams {
        cteam := CTeam{Team: t}
        user, err := t.Members()
        if err != nil {
            h.JSONR(c, http.StatusBadRequest, err)
            return
        }
        cteam.Users = user
        creatorName, err := t.GetCreatorName()
        if err != nil {
            log.Debug(err.Error())
        }
        cteam.TeamCreator = creatorName
        outputs = append(outputs, cteam)
    }
    h.JSONR(c, outputs)
    return
}

func GetTeamByName(c *gin.Context) {
    name := c.Params.ByName("name")
    if name == "" {
        h.JSONR(c, http.StatusBadRequest, "team name is missing")
        return
    }

    var team um.Team
    dt := db.Uic.Table("team").Where(&um.Team{Name: name}).Find(&team)                                  
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error) 
        return                                                                                           
    }   

    var uidarr []um.RelTeamUser
    dt = db.Uic.Table("rel_team_user").Select("uid").Where(&um.RelTeamUser{Tid: team.ID}).Find(&uidarr) 
    if dt.Error != nil {
        log.Debug(dt.Error)
    }

    var resp APIGetTeamOutput
    resp.Team = team
    resp.Users = []um.User{}                                                                            
    if len(uidarr) != 0 {
        uids := ""
        for idx, v := range uidarr {
            if idx == 0 {
                uids = fmt.Sprintf("%v", v.Uid)
            } else {
                uids = fmt.Sprintf("%v,%v", uids, v.Uid)
            }
        }
        log.Debugf("uids:%s", uids)

        var users []um.User
        db.Uic.Table("user").Where(fmt.Sprintf("id IN (%s)", uids)).Find(&users)
        resp.Users = users
    }
    h.JSONR(c, resp)
    return
}

func CreateTeam(c *gin.Context) {
    var cteam APICreateTeamInput
    err := c.Bind(&cteam)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)                                                                       
        return                                                                                           
    }

    // here creator is always 1--root
    team := um.Team{
        Name:    cteam.Name,
        Resume:  cteam.Resume,                                                                           
        Creator: 1,
    }

    dt := db.Uic.Table("team").Create(&team)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    var dt2 *gorm.DB
    if len(cteam.UserIDs) > 0 {
        for i := 0; i < len(cteam.UserIDs); i++ {
            dt2 = db.Uic.Create(&um.RelTeamUser{Tid: team.ID, Uid: cteam.UserIDs[i]})
            if dt2.Error != nil {
                err = dt2.Error
                break
            }
        }
        if err != nil {
            h.JSONR(c, http.StatusBadRequest, err)
            return
        }
    }
    h.JSONR(c, fmt.Sprintf("team created! Afftect row: %d, Affect refer: %d", dt.RowsAffected, len(cteam.UserIDs)))
    return
}

func UpdateTeam(c *gin.Context) {
    var cteam APIUpdateTeamInput
    err := c.Bind(&cteam)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    dt := db.Uic
    dt = dt.Table("team").Where("id = ?", cteam.ID)

    var team um.Team
    dt = dt.Find(&team)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }
    tm := um.Team{Name: cteam.Name, Resume: cteam.Resume}
    dt = db.Uic.Table("team").Where("id=?", cteam.ID).Update(&tm)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    err = bindUsers(db, cteam.ID, cteam.UserIDs)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
    } else {
        h.JSONR(c, "team updated!")
    }
}

func bindUsers(db model.DBPool, tid int, users []int) (err error) {
    var dt *gorm.DB
    uids, err := utils.ArrIntToString(users)
    if err != nil {
        return
    }

    //delete unbind users
    var needDeleteMan []um.RelTeamUser
    qPared := fmt.Sprintf("tid = %d AND NOT (uid IN (%v))", tid, uids)
    dt = db.Uic.Table("rel_team_user").Where(qPared).Find(&needDeleteMan)
    if dt.Error != nil {
        err = dt.Error
        return
    }
    if len(needDeleteMan) != 0 {
        for _, man := range needDeleteMan {
            dt = db.Uic.Delete(&man)
            if dt.Error != nil {
                err = dt.Error
                return
            }
        }
    }

    //insert bind users
    for _, i := range users {
        ur := um.RelTeamUser{Tid: int64(tid), Uid: int64(i)}
        db.Uic.Table("rel_team_user").Where(&ur).Find(&ur)
        if ur.ID == 0 {
            dt = db.Uic.Table("rel_team_user").Create(&ur)
        } else {
            //if record exist, do next
            continue
        }
        if dt.Error != nil {
            err = dt.Error
            return
        }
    }
    return
}

func DeleteTeam(c *gin.Context) {
    var err error
    teamIdStr := c.Params.ByName("id")
    teamIdTmp, err := strconv.Atoi(teamIdStr)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err.Error())
        return
    }
    teamId := int64(teamIdTmp)
    if teamId == 0 {
        h.JSONR(c, http.StatusBadRequest, "team_id is empty")
        return
    } else if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    dt := db.Uic.Table("team")
    dt = dt.Delete(&um.Team{ID: teamId})
    err = dt.Error

    var dt2 *gorm.DB
    if err != nil {
        h.JSONR(c, http.StatusExpectationFailed, err)
        return
    } else {
        dt2 = db.Uic.Where("tid = ?", teamId).Delete(um.RelTeamUser{})
    }
    h.JSONR(c, fmt.Sprintf("team %v is deleted. Affect row: %d / refer delete: %d", teamId, dt.RowsAffected, dt2.RowsAffected))
    return
}

func Routes(r *gin.Engine) {
    db = model.Con()

    userapi := r.Group("/v1/user")
    userapi.GET("", GetUsers)
    userapi.GET("/:name", GetUserByName)
    userapi.POST("", CreateUser)
    userapi.PUT("/profile", UpdateUserProfile)
    userapi.PUT("/password", UpdateUserPassword)
    userapi.GET("/:name/team", GetUserTeams)
    userapi.DELETE("", DeleteUser)

    teamapi := r.Group("/v1/team")
    teamapi.GET("", GetTeams)
    teamapi.GET("/:name", GetTeamByName)
    teamapi.POST("", CreateTeam)
    teamapi.PUT("", UpdateTeam)
    teamapi.DELETE("/:id", DeleteTeam)
}

