package alarm

import (
    "errors"
    "fmt"
    "net/http"
    "strings"
    "time"

    h "falcon-api/application/helper"
    "falcon-api/application/model"
    am "falcon-api/application/model/alarm"

    "github.com/gin-gonic/gin"
)

type APIGetAlarmListsInputs struct {
    StartTime     int64    `json:"startTime" form:"startTime"`
    EndTime       int64    `json:"endTime" form:"endTime"`
    Priority      int      `json:"priority" form:"priority"`
    Status        string   `json:"status" form:"status"`
    ProcessStatus string   `json:"process_status" form:"process_status"`
    Metrics       string   `json:"metrics" form:"metrics"`
    //id
    EventId       string   `json:"event_id" form:"event_id"`
    //number of reacord's limit on each page
    Limit         int      `json:"limit" form:"limit"`
    //pagging
    Page          int      `json:"page" form:"page"`
    //endpoints strategy template
    Endpoints     []string `json:"endpoints" form:"endpoints"`
    StrategyId    int      `json:"strategy_id" form:"strategy_id"`
    TemplateId    int      `json:"template_id" form:"template_id"`
}

type APIEventsGetInputs struct {
    StartTime int64  `json:"startTime" form:"startTime"`
    EndTime   int64  `json:"endTime" form:"endTime"`
    Status    int    `json:"status" form:"status" binding:"gte=-1,lte=1"`
    //event_caseId
    EventId   string `json:"event_id" form:"event_id" binding:"required"`
    //number of reacord's limit on each page
    Limit     int    `json:"limit" form:"limit"`
    //pagging
    Page      int    `json:"page" form:"page"`
}

type APIGetNotesOfAlarmInputs struct {
    StartTime int64  `json:"startTime" form:"startTime"`
    EndTime   int64  `json:"endTime" form:"endTime"`
    EventId   string `json:"event_id" form:"event_id"`
    Status    string `json:"status" form:"status"`
    Limit     int    `json:"limit" form:"limit"`
    Page      int    `json:"page" form:"page"`
}

type APIGetNotesOfAlarmOuput struct {
    EventCaseId string     `json:"event_caseId"`
    Note        string     `json:"note"`
    CaseId      string     `json:"case_id"`
    Status      string     `json:"status"`
    Timestamp   *time.Time `json:"timestamp"`
    UserName    string     `json:"user"`
}

var db model.DBPool

func (input APIGetAlarmListsInputs) checkInputsContain() error {
    if input.StartTime == 0 && input.EndTime == 0 {
        if input.EventId == "" && input.Endpoints == nil && input.StrategyId == 0 && input.TemplateId == 0 {
            return errors.New("startTime, endTime, event_id, endpoints, strategy_id or template_id, You have to at least pick one on the request.")
        }
    }
    return nil
}

func (s APIGetAlarmListsInputs) collectFilters() string {
    tmp := []string{}
    if s.StartTime != 0 {
        tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", s.StartTime))
    }
    if s.EndTime != 0 {
        tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", s.EndTime))
    }
    if s.Priority != -1 {
        tmp = append(tmp, fmt.Sprintf("priority = %d", s.Priority))
    }
    if s.Status != "" {
        status := ""
        statusTmp := strings.Split(s.Status, ",")
        for indx, n := range statusTmp {
            if indx == 0 {
                status = fmt.Sprintf(" status = '%s' ", n)
            } else {
                status = fmt.Sprintf(" %s OR status = '%s' ", status, n)
            }
        }
        status = fmt.Sprintf("( %s )", status)
        tmp = append(tmp, status)
    }
    if s.ProcessStatus != "" {
        pstatus := ""
        pstatusTmp := strings.Split(s.ProcessStatus, ",")
        for indx, n := range pstatusTmp {
            if indx == 0 {
                pstatus = fmt.Sprintf(" process_status = '%s' ", n)
            } else {
                pstatus = fmt.Sprintf(" %s OR process_status = '%s' ", pstatus, n)
            }
        }
        pstatus = fmt.Sprintf("( %s )", pstatus)
        tmp = append(tmp, pstatus)
    }
    if s.Metrics != "" {
        tmp = append(tmp, fmt.Sprintf("metrics regexp '%s'", s.Metrics))
    }
    if s.EventId != "" {
        tmp = append(tmp, fmt.Sprintf("id = '%s'", s.EventId))
    }
    if s.Endpoints != nil && len(s.Endpoints) != 0 {
        for i, ep := range s.Endpoints {
            s.Endpoints[i] = fmt.Sprintf("'%s'", ep)
        }
        tmp = append(tmp, fmt.Sprintf("endpoint in (%s)", strings.Join(s.Endpoints, ", ")))
    }
    if s.StrategyId != 0 {
        tmp = append(tmp, fmt.Sprintf("strategy_id = %d", s.StrategyId))
    }
    if s.TemplateId != 0 {
        tmp = append(tmp, fmt.Sprintf("template_id = %d", s.TemplateId))
    }
    filterStrTmp := strings.Join(tmp, " AND ")
    if filterStrTmp != "" {
        filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
    }
    return filterStrTmp
}

func (s APIEventsGetInputs) collectFilters() string {
    tmp := []string{}
    filterStrTmp := ""
    if s.StartTime != 0 { 
        tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", s.StartTime))
    }   
    if s.EndTime != 0 { 
        tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", s.EndTime))
    }   
    if s.EventId != "" {
        tmp = append(tmp, fmt.Sprintf("event_caseId = '%s'", s.EventId))
    }   
    if s.Status == 0 || s.Status == 1 { 
        tmp = append(tmp, fmt.Sprintf("status = %d", s.Status))
    }   
    if len(tmp) != 0 { 
        filterStrTmp = strings.Join(tmp, " AND ")
        filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
    }   
    return filterStrTmp
}

func (input APIGetNotesOfAlarmInputs) checkInputsContain() error {
    if input.StartTime == 0 && input.EndTime == 0 { 
        if input.EventId == "" {
            return errors.New("startTime, endTime OR event_id, You have to at least pick one on the request.")
        }   
    }   
    return nil 
}

func (s APIGetNotesOfAlarmInputs) collectFilters() string {
    tmp := []string{}
    if s.StartTime != 0 { 
        tmp = append(tmp, fmt.Sprintf("timestamp >= FROM_UNIXTIME(%v)", s.StartTime))
    }   
    if s.EndTime != 0 { 
        tmp = append(tmp, fmt.Sprintf("timestamp <= FROM_UNIXTIME(%v)", s.EndTime))
    }
    if s.Status != "" {
        tmp = append(tmp, fmt.Sprintf("status = '%s'", s.Status))
    }
    if s.EventId != "" {
        tmp = append(tmp, fmt.Sprintf("event_caseId = '%s'", s.EventId))
    }
    filterStrTmp := strings.Join(tmp, " AND ")
    if filterStrTmp != "" {
        filterStrTmp = fmt.Sprintf("WHERE %s", filterStrTmp)
    }
    return filterStrTmp
}


func GetAlarms(c *gin.Context) {
    var inputs APIGetAlarmListsInputs
    //set default
    inputs.Page = -1
    inputs.Priority = -1

    if err := c.Bind(&inputs); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "400",
        })
        return
    }   

    if err := inputs.checkInputsContain(); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "400",
        })
        return
    }   

    filterCollector := inputs.collectFilters()
    //for get correct table name
    f := am.EventCases{}
    cases := []am.EventCases{}
    perparedSql := ""
    //if no specific, will give return first 2000 records
    if inputs.Page == -1 {
        if inputs.Limit >= 2000 || inputs.Limit == 0 { 
            inputs.Limit = 2000
        }
        perparedSql = fmt.Sprintf("select * from %s %s order by timestamp DESC limit %d", f.TableName(), filterCollector, inputs.Limit)
    } else {
        //set the max limit of each page
        if inputs.Limit >= 50 {
            inputs.Limit = 50
        }
        perparedSql = fmt.Sprintf("select * from %s %s  order by timestamp DESC limit %d,%d", f.TableName(), filterCollector, inputs.Page, inputs.Limit)
    }
    db.Alarm.Raw(perparedSql).Find(&cases)
    h.JSONR(c, cases)
}

func GetEvents(c *gin.Context) {
    var inputs APIEventsGetInputs
    inputs.Status = -1
    if err := c.Bind(&inputs); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "400",
        })
        return
    }
    filterCollector := inputs.collectFilters()
    //for get correct table name
    f := am.Events{}
    evens := []am.Events{}
    if inputs.Limit == 0 || inputs.Limit >= 50 {
        inputs.Limit = 50
    }
    perparedSql := fmt.Sprintf("select id, event_caseId, cond, status, timestamp from %s %s order by timestamp DESC limit %d,%d", f.TableName(), filterCollector, inputs.Page, inputs.Limit)
    db.Alarm.Raw(perparedSql).Scan(&evens)
    h.JSONR(c, evens)
}

func GetNotesOfAlarm(c *gin.Context) {
    var inputs APIGetNotesOfAlarmInputs
    if err := c.Bind(&inputs); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "400",
        })
        return
    }

    if err := inputs.checkInputsContain(); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "400",
        })
        return
    }

    filterCollector := inputs.collectFilters()

    //for get correct table name
    f := am.EventNote{}
    notes := []am.EventNote{}
    if inputs.Limit == 0 || inputs.Limit >= 50 {
        inputs.Limit = 50
    }
    perparedSql := fmt.Sprintf(
        "select id, event_caseId, note, case_id, status, timestamp, user_id from %s %s order by timestamp DESC limit %d,%d",
        f.TableName(),
        filterCollector,
        inputs.Page,
        inputs.Limit,
    )
    db.Alarm.Raw(perparedSql).Scan(&notes)
    output := []APIGetNotesOfAlarmOuput{}
    for _, n := range notes {
        output = append(output, APIGetNotesOfAlarmOuput{
            EventCaseId: n.EventCaseId,
            Note:        n.Note,
            CaseId:      n.CaseId,
            Status:      n.Status,
            Timestamp:   n.Timestamp,
            UserName:    n.GetUserName(),
        })
    }
    h.JSONR(c, output)
}

func Routes(r *gin.Engine) {
    db = model.Con()

    alarmapi := r.Group("/v1/alarm")

    alarmapi.POST("/case", GetAlarms)
    alarmapi.GET("/event", GetEvents)
    alarmapi.GET("/note", GetNotesOfAlarm)
}

