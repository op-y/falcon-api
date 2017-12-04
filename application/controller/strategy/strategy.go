package strategy

import (
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
    "strconv"
    "strings"

    h "falcon-api/application/helper"
    "falcon-api/application/model"
    fpm "falcon-api/application/model/falcon_portal"

    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
)

type APICreateStrategyInput struct {
    Metric     string `json:"metric" binding:"required"`                                                 
    Tags       string `json:"tags"`                                                                      
    MaxStep    int    `json:"max_step" binding:"required"`                                               
    Priority   int    `json:"priority" binding:"exists"`                                                 
    Func       string `json:"func" binding:"required"`                                                   
    Op         string `json:"op" binding:"required"`                                                     
    RightValue string `json:"right_value" binding:"required"`                                            
    Note       string `json:"note"`                                                                      
    RunBegin   string `json:"run_begin"`                                                                 
    RunEnd     string `json:"run_end"`
    TplId      int64  `json:"tpl_id" binding:"required"`                                                 
}

type APIUpdateStrategyInput struct {                                                                     
    ID         int64  `json:"id" binding:"required"`                                                     
    Metric     string `json:"metric" binding:"required"`                                                 
    Tags       string `json:"tags"`                                                                      
    MaxStep    int    `json:"max_step" binding:"required"`                                               
    Priority   int    `json:"priority" binding:"exists"`                                                 
    Func       string `json:"func" binding:"required"`                                                   
    Op         string `json:"op" binding:"required"`                                                     
    RightValue string `json:"right_value" binding:"required"`                                            
    Note       string `json:"note"`
    RunBegin   string `json:"run_begin"`
    RunEnd     string `json:"run_end"`                                                                   
}

var db model.DBPool

func (this APICreateStrategyInput) CheckFormat() (err error) {                                           
    validOp := regexp.MustCompile(`^(>|=|<|!)(=)?$`)                                                     
    validRightValue := regexp.MustCompile(`^\-?\d+(\.\d+)?$`)                                            
    validTime := regexp.MustCompile(`^\d{2}:\d{2}$`)                                                     
    switch {
    case !validOp.MatchString(this.Op):
        err = errors.New("op's formating is not vaild")                                                  
    case !validRightValue.MatchString(this.RightValue):
        err = errors.New("right_value's formating is not vaild")
    case !validTime.MatchString(this.RunBegin) && this.RunBegin != "":
        err = errors.New("run_begin's formating is not vaild, please refer ex. 00:00")
    case !validTime.MatchString(this.RunEnd) && this.RunEnd != "":
        err = errors.New("run_end's formating is not vaild, please refer ex. 24:00")
    }
    return
}

func (this APIUpdateStrategyInput) CheckFormat() (err error) {
    validOp := regexp.MustCompile(`^(>|=|<|!)(=)?$`)
    validRightValue := regexp.MustCompile(`^\-?\d+(\.\d+)?$`)
    validTime := regexp.MustCompile(`^\d{2}:\d{2}$`)
    switch {
    case !validOp.MatchString(this.Op):
        err = errors.New("op's formating is not vaild")
    case !validRightValue.MatchString(this.RightValue):
        err = errors.New("right_value's formating is not vaild")
    case !validTime.MatchString(this.RunBegin) && this.RunBegin != "":
        err = errors.New("run_begin's formating is not vaild, please refer ex. 00:00")
    case !validTime.MatchString(this.RunEnd) && this.RunEnd != "":
        err = errors.New("run_end's formating is not vaild, please refer ex. 24:00")
    }
    return
}

func GetStrategys(c *gin.Context) {
    var strategys []fpm.Strategy

    tidtmp := c.DefaultQuery("tid", "") 
    if tidtmp == "" {
        h.JSONR(c, http.StatusBadRequest, "tid is missing")
        return
    }

    tid, err := strconv.Atoi(tidtmp)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    dt := db.FalconPortal.Where("tpl_id = ?", tid).Find(&strategys)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    h.JSONR(c, strategys)
    return
}

func GetStrategy(c *gin.Context) {
    sidtmp := c.Params.ByName("id")                                                                     
    if sidtmp == "" {                                                                                    
        h.JSONR(c, http.StatusBadRequest, "sid is missing")                                                          
        return
    }

    sid, err := strconv.Atoi(sidtmp)                                                                     
    if err != nil {                                                                                      
        h.JSONR(c, http.StatusBadRequest, err)                                                                       
        return                                                                                           
    }

    strategy := fpm.Strategy{ID: int64(sid)}
    if dt := db.FalconPortal.Find(&strategy); dt.Error != nil {                                                
        h.JSONR(c, http.StatusBadRequest, dt.Error)                                                                  
        return
    }

    h.JSONR(c, strategy)
    return
}

func CreateStrategy(c *gin.Context) {
    var inputs APICreateStrategyInput
    if err := c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }
    if err := inputs.CheckFormat(); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    strategy := fpm.Strategy{
        Metric:     inputs.Metric,
        Tags:       inputs.Tags,
        MaxStep:    inputs.MaxStep,
        Priority:   inputs.Priority,
        Func:       inputs.Func,
        Op:         inputs.Op,
        RightValue: inputs.RightValue,
        Note:       inputs.Note,
        RunBegin:   inputs.RunBegin,
        RunEnd:     inputs.RunEnd,
        TplId:      inputs.TplId,
    }
    dt := db.FalconPortal.Save(&strategy)
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }

    h.JSONR(c, "stragtegy created")
    return
}

func UpdateStrategy(c *gin.Context) {
    var inputs APIUpdateStrategyInput
    if err := c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }
    if err := inputs.CheckFormat(); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    strategy := fpm.Strategy{
        ID: inputs.ID,
    }
    if dt := db.FalconPortal.Find(&strategy); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("find strategy got error:%v", dt.Error))
        return
    }
    ustrategy := map[string]interface{}{
        "Metric":     inputs.Metric,
        "Tags":       inputs.Tags,
        "MaxStep":    inputs.MaxStep,
        "Priority":   inputs.Priority,
        "Func":       inputs.Func,
        "Op":         inputs.Op,
        "RightValue": inputs.RightValue,
        "Note":       inputs.Note,
        "RunBegin":   inputs.RunBegin,
        "RunEnd":     inputs.RunEnd}
    if dt := db.FalconPortal.Model(&strategy).Where("id = ?", strategy.ID).Update(ustrategy); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }

    h.JSONR(c, fmt.Sprintf("stragtegy:%d has been updated", strategy.ID))
    return
}

func DeleteStrategy(c *gin.Context) {
    sidtmp := c.Params.ByName("id")
    if sidtmp == "" {
        h.JSONR(c, http.StatusBadRequest, "sid is missing")
        return
    }
    sid, err := strconv.Atoi(sidtmp)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    strategy := fpm.Strategy{ID: int64(sid)}
    if dt := db.FalconPortal.Delete(&strategy); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    h.JSONR(c, fmt.Sprintf("strategy:%d has been deleted", sid))
    return
}

func QueryMetrics(c *gin.Context) {
    //filePath := "./data/metric"
    filePath := viper.GetString("metric_list_file")
    if filePath == "" {
        filePath = "./data/metric"
    }
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }
    metrics := strings.Split(string(data), "\n")
    h.JSONR(c, metrics)
    return
}

func Routes(r *gin.Engine) {
    db = model.Con()

    strategyapi := r.Group("/v1/strategy")

    strategyapi.GET("", GetStrategys)
    strategyapi.GET("/:id", GetStrategy)
    strategyapi.POST("", CreateStrategy)
    strategyapi.PUT("", UpdateStrategy)
    strategyapi.DELETE("/:id", DeleteStrategy)

    metricapi := r.Group("/v1/metric")
    metricapi.GET("", QueryMetrics)
}

