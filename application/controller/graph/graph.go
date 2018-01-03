package graph

import (
    "errors"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"

    h "falcon-api/application/helper"
    "falcon-api/application/model"
    gclient "falcon-api/application/graph"
    gm "falcon-api/application/model/graph"
    //"falcon-api/application/utils"

    cmodel "github.com/open-falcon/falcon-plus/common/model"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    log "github.com/Sirupsen/logrus"
    tcache "github.com/toolkits/cache/localcache/timedcache"
)

type APIEndpointObjGetInputs struct {
    Endpoints []string `json:"endpoints" form:"endpoints"`
    Deadline  int64    `json:"deadline" form:"deadline"`
}

type APIEndpointRegexpQueryInputs struct {
    Q     string `json:"q" form:"q"`
    Label string `json:"tags" form:"tags"`                                                               
    Limit int    `json:"limit" form:"limit"`                                                             
    Page  int    `json:"page" form:"page"`                                                               
}

type APIQueryGraphDrawData struct {
    HostNames []string `json:"hostnames" binding:"required"`                                             
    Counters  []string `json:"counters" binding:"required"`                                              
    ConsolFun string   `json:"consol_fun" binding:"required"`                                            
    StartTime int64    `json:"start_time" binding:"required"`                                            
    EndTime   int64    `json:"end_time" binding:"required"`                                              
    Step      int      `json:"step"`
}

type APIGraphDeleteCounterInputs struct {
    Endpoints []string `json:"endpoints" binding:"required"`
    Counters  []string `json:"counters" binding:"required"`
}

var localStepCache = tcache.New(600*time.Second, 60*time.Second)
var db model.DBPool

func GetEndpointObject(c *gin.Context) {
    inputs := APIEndpointObjGetInputs{
        Deadline: 0,
    }   
    if err := c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    if len(inputs.Endpoints) == 0 { 
        h.JSONR(c, http.StatusBadRequest, "endpoints missing")
        return
    }

    var result []gm.Endpoint = []gm.Endpoint{}
    dt := db.Graph.Table("endpoint").
        Where("endpoint in (?) and ts >= ?", inputs.Endpoints, inputs.Deadline).
        Scan(&result)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    endpoints := []map[string]interface{}{}
    for _, r := range result {
        endpoints = append(endpoints, map[string]interface{}{"id": r.ID, "endpoint": r.Endpoint, "ts": r.Ts})
    }

    h.JSONR(c, endpoints)
}

func GetEndpointByRegExp(c *gin.Context) {
    inputs := APIEndpointRegexpQueryInputs{
        //set default is 500
        Limit: 500,
        Page:  1,
    }
    if err := c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)                                                                       
        return
    }
    if inputs.Q == "" && inputs.Label == "" {
        h.JSONR(c, http.StatusBadRequest, "q and labels are all missing")                                
        return
    }                                                                                                    

    labels := []string{}
    if inputs.Label != "" {
        labels = strings.Split(inputs.Label, ",")
    }
    qs := []string{}
    if inputs.Q != "" {
        qs = strings.Split(inputs.Q, " ")
    }

    var offset int = 0
    if inputs.Page > 1 {
        offset = (inputs.Page - 1) * inputs.Limit
    }

    var endpoint []gm.Endpoint
    var endpoint_id []int
    var dt *gorm.DB
    if len(labels) != 0 {
        dt = db.Graph.Table("endpoint_counter").Select("distinct endpoint_id")
        for _, trem := range labels {
            dt = dt.Where(" counter like ? ", "%"+strings.TrimSpace(trem)+"%")
        }
        dt = dt.Limit(inputs.Limit).Offset(offset).Pluck("distinct endpoint_id", &endpoint_id)
        if dt.Error != nil {
            h.JSONR(c, http.StatusBadRequest, dt.Error)
            return
        }
    }
    if len(qs) != 0 {
        dt = db.Graph.Table("endpoint").
            Select("endpoint, id")
        if len(endpoint_id) != 0 {
            dt = dt.Where("id in (?)", endpoint_id)
        }

        for _, trem := range qs {
            dt = dt.Where(" endpoint regexp ? ", strings.TrimSpace(trem))
        }
        dt.Limit(inputs.Limit).Offset(offset).Scan(&endpoint)
    } else if len(endpoint_id) != 0 {
        dt = db.Graph.Table("endpoint").
            Select("endpoint, id").
            Where("id in (?)", endpoint_id).
            Scan(&endpoint)
    }
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    endpoints := []map[string]interface{}{}
    for _, e := range endpoint {
        endpoints = append(endpoints, map[string]interface{}{"id": e.ID, "endpoint": e.Endpoint})
    }

    h.JSONR(c, endpoints)
}

func GetEndpointCounterByRegExp(c *gin.Context) {
    hostQuery := c.DefaultQuery("hostQuery", "")
    if hostQuery == "" {
        h.JSONR(c, http.StatusBadRequest, "hostQuery is missing")
        return
    }

    hosts := []string{}
    hosts = strings.Split(hostQuery, ",")

    metricQuery := c.DefaultQuery("metricQuery", ".+")
    
    limitTmp := c.DefaultQuery("limit", "500")
    limit, err := strconv.Atoi(limitTmp)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    pageTmp := c.DefaultQuery("page", "1")                                                               
    page, err := strconv.Atoi(pageTmp)                                                                   
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    var offset int = 0                                                                                   
    if page > 1 {                                                                                        
        offset = (page - 1) * limit
    }

    // query endpoint ids by hostname
    var endpoint_id []int
    var dt *gorm.DB
    dt = db.Graph.Table("endpoint").Select("id")
    for _, host := range hosts {
        //dt = dt.Where(" endpoint regexp ? ", strings.TrimSpace(host))
        dt = dt.Where(" endpoint like ? ", "%"+strings.TrimSpace(host)+"%")
    }
    dt = dt.Pluck("id", &endpoint_id)

    // prepare condition string
    condition := ""
    if len(endpoint_id) == 0 {
        h.JSONR(c, http.StatusBadRequest, "no endpoint id, please check your input info.")
        return
    } else {
	    for idx, id := range endpoint_id {
			if idx == 0 {
				condition = fmt.Sprintf("%d", id)
			} else {
				condition = fmt.Sprintf("%s, %d", condition, id)
			}
        }
    }

    // query counters
    var counters []gm.EndpointCounter
    dt = db.Graph.Table("endpoint_counter").Select("endpoint_id, counter, step, type").Where(fmt.Sprintf("endpoint_id IN %s", condition))
    if metricQuery != "" {
        qs := strings.Split(metricQuery, " ")
        if len(qs) > 0 {
            for _, term := range qs {
                dt = dt.Where("counter regexp ?", strings.TrimSpace(term))
            }
        }
    }
    dt = dt.Limit(limit).Offset(offset).Scan(&counters)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    // prepare response
    countersResp := []interface{}{}
    for _, c := range counters {
        countersResp = append(countersResp, map[string]interface{}{
            "endpoint_id": c.EndpointID,
            "counter":     c.Counter,
            "step":        c.Step,
            "type":        c.Type,
        })
    }
    h.JSONR(c, countersResp)
    return
}

func GetGraphDrawData(c *gin.Context) {
    var inputs APIQueryGraphDrawData
    var err error
    if err = c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return                                                                                           
    }
    respData := []*cmodel.GraphQueryResponse{}
    for _, host := range inputs.HostNames {
        for _, counter := range inputs.Counters {
            var step int
            if inputs.Step > 0 {                                                                         
                step = inputs.Step                                                                       
            } else {
                step, err = getCounterStep(host, counter)
                if err != nil {
                    continue
                }
            }
            data, _ := fetchData(host, counter, inputs.ConsolFun, inputs.StartTime, inputs.EndTime, step)
            respData = append(respData, data)
        }
    }
    h.JSONR(c, respData)
}

func GetGraphLastPoint(c *gin.Context) {
    var inputs []cmodel.GraphLastParam
    if err := c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }   
    respData := []*cmodel.GraphLastResp{}

    for _, param := range inputs {
        one_resp, err := gclient.Last(param)
        if err != nil {
            log.Warn("query last point from graph fail:", err)
        } else {
            respData = append(respData, one_resp)
        }   
    }   

    h.JSONR(c, respData)
}

func getCounterStep(endpoint, counter string) (step int, err error) {
    cache_key := fmt.Sprintf("step:%s/%s", endpoint, counter)
    s, found := localStepCache.Get(cache_key)
    if found && s != nil {
        step = s.(int)
        return
    }

    var rows []int
    dt := db.Graph.Raw(`select a.step from endpoint_counter as a, endpoint as b
        where b.endpoint = ? and a.endpoint_id = b.id and a.counter = ? limit 1`, endpoint, counter).Scan(&rows)
    if dt.Error != nil {
        err = dt.Error
        return
    }
    if len(rows) == 0 {
        err = errors.New("empty result")
        return
    }
    step = rows[0]
    localStepCache.Set(cache_key, step, tcache.DefaultExpiration)

    return
}

func fetchData(hostname string, counter string, consolFun string, startTime int64, endTime int64, step int) (resp *cmodel.GraphQueryResponse, err error) {
    qparm := gclient.GenQParam(hostname, counter, consolFun, startTime, endTime, step)
    resp, err = gclient.QueryOne(qparm)
    if err != nil {
        log.Debugf("query graph got error: %s", err.Error())
    }
    return
}

func DeleteGraphEndpoint(c *gin.Context) {
    var inputs []string = []string{}
    if err := c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)                                                                       
        return
    }   

    type DBRows struct {
        Endpoint  string
        CounterId int                                                                                    
        Counter   string                                                                                 
        Type      string
        Step      int
    }

    rows := []DBRows{}
    dt := db.Graph.Raw(
        `select a.endpoint, b.id AS counter_id, b.counter, b.type, b.step from endpoint as a, endpoint_counter as b
        where b.endpoint_id = a.id
        AND a.endpoint in (?)`, inputs).Scan(&rows)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    var affected_counter int64 = 0
    var affected_endpoint int64 = 0

    if len(rows) > 0 {
        var params []*cmodel.GraphDeleteParam = []*cmodel.GraphDeleteParam{}
        for _, row := range rows {
            param := &cmodel.GraphDeleteParam{
                Endpoint: row.Endpoint,
                DsType:   row.Type,
                Step:     row.Step,
            }
            fields := strings.SplitN(row.Counter, "/", 2)
            if len(fields) == 1 {
                param.Metric = fields[0]
            } else if len(fields) == 2 {
                param.Metric = fields[0]
                param.Tags = fields[1]
            } else {
                log.Error("invalid counter", row.Counter)
                continue
            }
            params = append(params, param)
        }
        gclient.Delete(params)
    }

    tx := db.Graph.Begin()

    if len(rows) > 0 {
        var cids []int = make([]int, len(rows))
        for i, row := range rows {
            cids[i] = row.CounterId
        }

        dt = tx.Table("endpoint_counter").Where("id in (?)", cids).Delete(&gm.EndpointCounter{})
        if dt.Error != nil {
            h.JSONR(c, http.StatusBadRequest, dt.Error)
            tx.Rollback()
            return
        }
        affected_counter = dt.RowsAffected

        dt = tx.Exec(`delete from tag_endpoint where endpoint_id in
            (select id from endpoint where endpoint in (?))`, inputs)
        if dt.Error != nil {
            h.JSONR(c, http.StatusBadRequest, dt.Error)
            tx.Rollback()
            return
        }
    }

    dt = tx.Table("endpoint").Where("endpoint in (?)", inputs).Delete(&gm.Endpoint{})
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        tx.Rollback()
        return
    }
    affected_endpoint = dt.RowsAffected

    tx.Commit()

    h.JSONR(c, map[string]int64{
        "affected_endpoint": affected_endpoint,
        "affected_counter":  affected_counter,
    })
}

func DeleteGraphCounter(c *gin.Context) {
    var inputs APIGraphDeleteCounterInputs = APIGraphDeleteCounterInputs{}
    if err := c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    type DBRows struct {
        Endpoint  string
        CounterId int
        Counter   string
        Type      string
        Step      int
    }

    rows := []DBRows{}
    dt := db.Graph.Raw(`select a.endpoint, b.id AS counter_id, b.counter, b.type, b.step from endpoint as a,
        endpoint_counter as b
        where b.endpoint_id = a.id
        AND a.endpoint in (?)
        AND b.counter in (?)`, inputs.Endpoints, inputs.Counters).Scan(&rows)
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }
    if len(rows) == 0 {
        h.JSONR(c, map[string]int64{
            "affected_counter": 0,
        })
        return
    }

    var params []*cmodel.GraphDeleteParam = []*cmodel.GraphDeleteParam{}
    for _, row := range rows {
        param := &cmodel.GraphDeleteParam{
            Endpoint: row.Endpoint,
            DsType:   row.Type,
            Step:     row.Step,
        }
        fields := strings.SplitN(row.Counter, "/", 2)
        if len(fields) == 1 {
            param.Metric = fields[0]
        } else if len(fields) == 2 {
            param.Metric = fields[0]
            param.Tags = fields[1]
        } else {
            log.Error("invalid counter", row.Counter)
            continue
        }
        params = append(params, param)
    }
    gclient.Delete(params)

    tx := db.Graph.Begin()
    var cids []int = make([]int, len(rows))
    for i, row := range rows {
        cids[i] = row.CounterId
    }

    dt = tx.Table("endpoint_counter").Where("id in (?)", cids).Delete(&gm.EndpointCounter{})
    if dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        tx.Rollback()
        return
    }
    affected_counter := dt.RowsAffected
    tx.Commit()

    h.JSONR(c, map[string]int64{
        "affected_counter": affected_counter,
    })
}

func Routes(r *gin.Engine) {
    db = model.Con()

    graphapi := r.Group("/v1/graph")

    graphapi.GET("/endpoint-object", GetEndpointObject)
    graphapi.GET("/endpoint", GetEndpointByRegExp)
    graphapi.GET("/endpoint-counter", GetEndpointCounterByRegExp)
    graphapi.POST("/history", GetGraphDrawData)
    graphapi.POST("/last-point", GetGraphLastPoint)
    graphapi.DELETE("/endpoint", DeleteGraphEndpoint)
    graphapi.DELETE("/counter", DeleteGraphCounter)
}

