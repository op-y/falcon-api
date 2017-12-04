package template

import (
    "fmt"
    "net/http"
    "strconv"

    h "falcon-api/application/helper"
    "falcon-api/application/model"
    fpm "falcon-api/application/model/falcon_portal"

    log "github.com/Sirupsen/logrus"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
)

type CTemplate struct {
    Template   fpm.Template `json:"template"`
    ParentName string       `json:"parent_name"`
}

type APIGetTemplatesOutput struct {
    Templates []CTemplate `json:"templates"`
}

type APICreateTemplateInput struct {
    Name     string `json:"name" binding:"required"`
    ParentID int64  `json:"parent_id" binding:"exists"`
    ActionID int64  `json:"action_id"`
}

type APIUpdateTemplateInput struct {
    Name     string `json:"name" binding:"required"`                                                     
    ParentID int64  `json:"parent_id" binding:"exists"`                                                  
    TplID    int64  `json:"tpl_id" binding:"required"`                                                   
}

type APICreateActionToTmplateInput struct {                                                              
    UIC                string `json:"uic" binding:"exists"`                                              
    URL                string `json:"url" binding:"exists"`                                              
    Callback           int    `json:"callback" binding:"exists"`                                         
    BeforeCallbackSMS  int    `json:"before_callback_sms" binding:"exists"`                              
    AfterCallbackSMS   int    `json:"after_callback_sms" binding:"exists"`                               
    BeforeCallbackMail int    `json:"before_callback_mail" binding:"exists"`                             
    AfterCallbackMail  int    `json:"after_callback_mail" binding:"exists"`                              
    TplId              int64  `json:"tpl_id" binding:"required"`                                         
}

type APIUpdateActionToTmplateInput struct {
    ID                 int64  `json:"id" binding:"required"`
    UIC                string `json:"uic" binding:"exists"`
    URL                string `json:"url" binding:"exists"`
    Callback           int    `json:"callback" binding:"exists"`
    BeforeCallbackSMS  int    `json:"before_callback_sms" binding:"exists"`
    AfterCallbackSMS   int    `json:"after_callback_sms" binding:"exists"`
    BeforeCallbackMail int    `json:"before_callback_mail" binding:"exists"`
    AfterCallbackMail  int    `json:"after_callback_mail" binding:"exists"`
}

var db model.DBPool

func GetTemplates(c *gin.Context) {
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

    var dt *gorm.DB
    var templates []fpm.Template
    q := c.DefaultQuery("q", ".+")
    if limit != -1 && page != -1 {
        dt = db.FalconPortal.Raw(
            fmt.Sprintf("SELECT * from tpl WHERE tpl_name regexp %s limit %d,%d", q, page, limit)).Scan(&templates)
    } else {
        dt = db.FalconPortal.Where("tpl_name regexp ?", q).Find(&templates)
    }

    if dt.Error != nil {
        log.Infof(dt.Error.Error())
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    output := APIGetTemplatesOutput{}
    output.Templates = []CTemplate{}
    for _, t := range templates {
        var pname string
        pname, err := t.FindParentName()
        if err != nil {
            h.JSONR(c, http.StatusBadRequest, err)
            return
        }
        output.Templates = append(output.Templates, CTemplate{
            Template:   t,
            ParentName: pname,
        })
    }
    h.JSONR(c, output)
    return
}

func GetTemplate(c *gin.Context) {
    tplidtmp := c.Params.ByName("id")                                                                
    if tplidtmp == "" {
        h.JSONR(c, http.StatusBadRequest, "tpl_id is missing")                                                       
        return
    }

    tplId, err := strconv.Atoi(tplidtmp)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    var tpl fpm.Template
    if dt := db.FalconPortal.Find(&tpl, tplId); dt.Error != nil {                                              
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }                                                                 

    var stratges []fpm.Strategy
    dt := db.FalconPortal.Where("tpl_id = ?", tplId).Find(&stratges)                                           
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return                                                                  
    }

    action := fpm.Action{}
    if tpl.ActionID != 0 {
        if dt = db.FalconPortal.Find(&action, tpl.ActionID); dt.Error != nil {
            h.JSONR(c, http.StatusBadRequest, dt.Error)
            return
        }
    }

    pname, _ := tpl.FindParentName()
    h.JSONR(c, map[string]interface{}{
        "template":    tpl,
        "stratges":    stratges,
        "action":      action,
        "parent_name": pname,
    })
    return
}

func CreateTemplate(c *gin.Context) {
    var inputs APICreateTemplateInput
    err := c.Bind(&inputs)
    log.Debugf("CreateTemplate input: %v", inputs)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    if inputs.Name == "" {
        h.JSONR(c, http.StatusBadRequest, "input name is empty, please check it")
        return
    }

    template := fpm.Template{
        Name:       inputs.Name,
        ParentID:   inputs.ParentID,
        ActionID:   inputs.ActionID,
        CreateUser: "root",
    }

    dt := db.FalconPortal.Table("tpl").Save(&template)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    h.JSONR(c, "template created")
    return
}

func UpdateTemplate(c *gin.Context) {
    var inputs APIUpdateTemplateInput
    err := c.Bind(&inputs)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    var tpl fpm.Template
    if dt := db.FalconPortal.Find(&tpl, inputs.TplID); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    utpl := map[string]interface{}{
        "Name":     inputs.Name,
        "ParentID": inputs.ParentID,
    }

    if dt := db.FalconPortal.Model(&tpl).Where("id = ?", inputs.TplID).Update(utpl); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    h.JSONR(c, "template updated")
    return
}

func DeleteTemplate(c *gin.Context) {
    tidTmp, _ := c.Params.Get("id")
    if tidTmp == "" {
        h.JSONR(c, http.StatusBadRequest, "template id is missing")
        return
    }

    tplId, err := strconv.Atoi(tidTmp)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    tx := db.FalconPortal.Begin()

    var tpl fpm.Template
    if dt := tx.Find(&tpl, tplId); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        tx.Rollback()
        return
    }

    //delete template
    actionId := tpl.ActionID
    if dt := tx.Delete(&tpl); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        tx.Rollback()
        return
    }

    //delete action
    if actionId != 0 {
        if dt := tx.Delete(&fpm.Action{}, actionId); dt.Error != nil {
            h.JSONR(c, http.StatusBadRequest, dt.Error)
            tx.Rollback()
            return
        }
    }

    //delete strategy
    if dt := tx.Where("tpl_id = ?", tplId).Delete(&fpm.Strategy{}); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        tx.Rollback()
        return
    }

    //delete grp_tpl
    if dt := tx.Where("tpl_id = ?", tplId).Delete(&fpm.GrpTpl{}); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        tx.Rollback()
        return
    }

    tx.Commit()
    h.JSONR(c, fmt.Sprintf("template %d has been deleted", tplId))
    return
}

func CreateActionToTmplate(c *gin.Context) {
    var inputs APICreateActionToTmplateInput
    err := c.Bind(&inputs)
    if err != nil {                                                           
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }                                                                 
    action := fpm.Action{
        UIC:                inputs.UIC,
        URL:                inputs.URL,                                                                  
        Callback:           inputs.Callback,                                                             
        BeforeCallbackSMS:  inputs.BeforeCallbackSMS,                                                    
        BeforeCallbackMail: inputs.BeforeCallbackMail,
        AfterCallbackMail:  inputs.AfterCallbackMail,
        AfterCallbackSMS:   inputs.AfterCallbackSMS,
    }

    tx := db.FalconPortal.Begin()

    if dt := tx.Table("action").Save(&action); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        tx.Rollback()
        return
    }

    var lid []int
    tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &lid)
    aid := lid[0]
    var tpl fpm.Template
    if dt := tx.Find(&tpl, inputs.TplId); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, fmt.Sprintf("template: %d ; %s", inputs.TplId, dt.Error.Error()))
        tx.Rollback()
        return
    }

    dt := tx.Model(&tpl).UpdateColumns(fpm.Template{ActionID: int64(aid)})
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        tx.Rollback()
        return
    }

    tx.Commit()

    h.JSONR(c, fmt.Sprintf("action is created and bind to template: %d", inputs.TplId))
    return
}

func UpdateActionToTmplate(c *gin.Context) {
    var inputs APIUpdateActionToTmplateInput
    err := c.BindJSON(&inputs)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    var action fpm.Action

    tx := db.FalconPortal.Begin()

    if dt := tx.Find(&action, inputs.ID); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        tx.Rollback()
        return
    }

    uaction := map[string]interface{}{
        "UIC":                inputs.UIC,
        "URL":                inputs.URL,
        "Callback":           inputs.Callback,
        "BeforeCallbackSMS":  inputs.BeforeCallbackSMS,
        "BeforeCallbackMail": inputs.BeforeCallbackMail,
        "AfterCallbackMail":  inputs.AfterCallbackMail,
        "AfterCallbackSMS":   inputs.AfterCallbackSMS,
    }
    dt := tx.Model(&action).Where("id = ?", inputs.ID).Update(uaction)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        tx.Rollback()
        return
    }

    tx.Commit()

    h.JSONR(c, fmt.Sprintf("action is updated, row affected: %d", dt.RowsAffected))
    return
}

func GetActionByID(c *gin.Context) {
    aid := c.Param("id") 
    act_id, err := strconv.Atoi(aid)                                                                     
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, "invalid action id")                                                       
        return
    }   

    act := fpm.Action{}
    dt := db.FalconPortal.Table("action").Where("id = ?", act_id).First(&act)                                  
    if dt.Error != nil {                                                                          
        h.JSONR(c, http.StatusBadRequest, dt.Error)                                                                  
        return
    }   

    h.JSONR(c, act) 
}

func Routes(r *gin.Engine) {
    db = model.Con()

    tplapi := r.Group("/v1/template")
    tplapi.GET("", GetTemplates)
    tplapi.GET("/:id", GetTemplate)
    tplapi.POST("", CreateTemplate)
    tplapi.PUT("", UpdateTemplate)
    tplapi.DELETE("/:id", DeleteTemplate)

    tplapi.POST("/action", CreateActionToTmplate)
    tplapi.PUT("/action", UpdateActionToTmplate)

    actionapi := r.Group("/v1/action")
    actionapi.GET("/:id", GetActionByID)
}

