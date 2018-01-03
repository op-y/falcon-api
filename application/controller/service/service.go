package service

import (
    "errors"
    "fmt"
    "net/http"
    "regexp"
    "strconv"

    h "falcon-api/application/helper"
    "falcon-api/application/model"
    fpm "falcon-api/application/model/falcon_portal"
    u "falcon-api/application/utils"

    log "github.com/Sirupsen/logrus"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
)

type APICrateHostGroup struct {
    Name string `json:"name" binding:"required"`                                                         
}

type APIBindHostToHostGroupInput struct {                                                                
    Hosts         []string `json:"hosts" binding:"required"`
    HostGroupName string   `json:"hostgroup_name" binding:"required"`                                        
}

type APIUnBindAHostToHostGroup struct {                                                                  
    HostName      string `json:"hosts" binding:"required"`                                                
    HostGroupName string `json:"hostgroup_name" binding:"required"`                                           
}

type APIHostGroupInputs struct {
    OldName string `json:"old_name" binding:"required"`
    NewName string `json:"new_name" binding:"required"`
}

type APIPatchHostGroupHost struct {
    Action string   `json:"action" binding:"required"`
    Hosts  []string `json:"hosts" binding:"required"`
}

type APICreateAggregatorInput struct {                                                                   
    GroupName   string `json:"hostgroup_name" binding:"required"`                                          
    Numerator   string `json:"numerator" binding:"required"`                                             
    Denominator string `json:"denominator" binding:"required"`                                           
    Endpoint    string `json:"endpoint" binding:"required"`                                              
    Metric      string `json:"metric" binding:"required"`                                                
    Tags        string `json:"tags" binding:"exists"`
    Step        int    `json:"step" binding:"required"`
}

type APIUpdateAggregatorInput struct {
    ID          int64  `json:"id" binding:"required"`                                                    
    Numerator   string `json:"numerator" binding:"required"`                                             
    Denominator string `json:"denominator" binding:"required"`
    Endpoint    string `json:"endpoint" binding:"required"`                                              
    Metric      string `json:"metric" binding:"required"`                                                
    Tags        string `json:"tags" binding:"exists"`                                                    
    Step        int    `json:"step" binding:"required"`                                                  
}

type APIBindTemplateToGroupInputs struct {
    TplID     int64  `json:"tpl_id"`                                                                          
    GroupName string `json:"grp_name"`                                                                          
}

type APIUnBindTemplateToGroupInputs struct {
    TplID     int64  `json:"tpl_id"`
    GroupName string `json:"grp_name"`
}

var db model.DBPool

func GetServices(c *gin.Context) {
    var (
        limit int 
        page  int 
        err   error
    )
    pageTmp := c.DefaultQuery("page", "") 
    limitTmp := c.DefaultQuery("limit", "") 
    q := c.DefaultQuery("q", ".+")
    page, limit, err = h.PageParser(pageTmp, limitTmp)
    if err != nil {
        h.JSONR(c, http.StatusBadRequest, err.Error())
        return
    }

    var hostgroups []fpm.HostGroup
    var dt *gorm.DB
    if limit != -1 && page != -1 {
        dt = db.FalconPortal.Raw(fmt.Sprintf("SELECT * from grp  where grp_name regexp '%s' limit %d,%d", q, page, limit)).Scan(&hostgroups)
    } else {
        dt = db.FalconPortal.Table("grp").Where("grp_name regexp ?", q).Find(&hostgroups)
    }   
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }   
    h.JSONR(c, hostgroups)
    return
}

func GetServiceByName(c *gin.Context) {
    grpName := c.Params.ByName("name")
    q := c.DefaultQuery("q", ".+") 
    if grpName == "" {
        h.JSONR(c, http.StatusBadRequest, "service name is missing")                                                       
        return
    }

    hostgroup := fpm.HostGroup{Name: grpName}
    if dt := db.FalconPortal.Where("grp_name = ?", hostgroup.Name).First(&hostgroup); dt.Error != nil {                                               
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)                                                                
        return
    }

    hosts := []fpm.Host{}
    grpHosts := []fpm.GrpHost{} 
    if dt := db.FalconPortal.Where("grp_id = ?", hostgroup.ID).Find(&grpHosts); dt.Error != nil {                     
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }
    for _, grph := range grpHosts {
        var host fpm.Host
        db.FalconPortal.Find(&host, grph.HostID)                                                               
        if host.ID != 0 {
            if ok, err := regexp.MatchString(q, host.Hostname); ok == true && err == nil {               
                hosts = append(hosts, host)
            }
        }                                                                                                
    }   
    h.JSONR(c, map[string]interface{}{
        "hostgroup": hostgroup,
        "hosts":     hosts,                                                                              
    })  
    return
}

func CreateService(c *gin.Context) {
    var inputs APICrateHostGroup
    if err := c.Bind(&inputs); err != nil { 
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    // service(unit) create by root
    hostgroup := fpm.HostGroup{Name: inputs.Name, CreateUser: "root", ComeFrom: 1}                      
    if dt := db.FalconPortal.Create(&hostgroup); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }                                                                                                    
    h.JSONR(c, hostgroup)                                                                                
    return
}

func BindInstanceToService(c *gin.Context) {
    var inputs APIBindHostToHostGroupInput                                                               
    if err := c.Bind(&inputs); err != nil {                                                              
        h.JSONR(c, http.StatusBadRequest, err)
        return                                                                                         
    }

    hostgroup := fpm.HostGroup{Name: inputs.HostGroupName}                                                     
    if dt := db.FalconPortal.Where("grp_name = ?", hostgroup.Name).First(&hostgroup); dt.Error != nil {                                               
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }

    tx := db.FalconPortal.Begin()
    if dt := tx.Where("grp_id = ?", hostgroup.ID).Delete(&fpm.GrpHost{}); dt.Error != nil {                
        h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("delete grp_host got error: %v", dt.Error))                  
        dt.Rollback()                                                                                    
        return
    }   

    var ids []int64
    for _, host := range inputs.Hosts {                                                                  
        ahost := fpm.Host{Hostname: host}                                                                  
        var id int64
        var ok bool
        if id, ok = ahost.Existing(); ok {
            ids = append(ids, id)
        } else {
            if dt := tx.Save(&ahost); dt.Error != nil {                                                  
                h.JSONR(c, http.StatusExpectationFailed, dt.Error)                                                        
                tx.Rollback()
                return
            }
            id = ahost.ID
            ids = append(ids, id)
        }
        if dt := tx.Debug().Create(&fpm.GrpHost{GrpID: hostgroup.ID, HostID: id}); dt.Error != nil {
            h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("create grphost got error: %s , grp_id: %v, host_id: %v", dt.Error, hostgroup.ID, id))
            tx.Rollback()
            return
        }
    }
    tx.Commit()

    h.JSONR(c, fmt.Sprintf("%v bind to hostgroup: %v", ids, hostgroup.ID))
    return
}

func UnbindAInstanceToService(c *gin.Context) {
    var inputs APIUnBindAHostToHostGroup                                                                 
    if err := c.Bind(&inputs); err != nil {                                                              
        h.JSONR(c, http.StatusBadRequest, err)                                                                       
        return
    }   

    hostgroup := fpm.HostGroup{Name: inputs.HostGroupName}                                                     
    if dt := db.FalconPortal.Where("grp_name = ?", hostgroup.Name).First(&hostgroup); dt.Error != nil {                                           
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    host := fpm.Host{Hostname: inputs.HostName}                                                     
    if dt := db.FalconPortal.Where("hostname = ?", host.Hostname).First(&host); dt.Error != nil {                                           
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    if dt := db.FalconPortal.Where("grp_id = ? AND host_id = ?", hostgroup.ID, host.ID).Delete(&fpm.GrpHost{}); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return                                                                                           
    }

    h.JSONR(c, fmt.Sprintf("unbind host:%s of hostgroup: %s", inputs.HostName, inputs.HostGroupName))        
    return
}

func UpdateService(c *gin.Context) {
    var inputs APIHostGroupInputs
    err := c.BindJSON(&inputs)
    switch {
    case err != nil:
        h.JSONR(c, http.StatusBadRequest, err)
        return
    case u.HasDangerousCharacters(inputs.NewName):                                                          
        h.JSONR(c, http.StatusBadRequest, "new_name is invalid")                                                     
        return
    }   

    hostgroup := fpm.HostGroup{Name: inputs.OldName}                                                           
    if dt := db.FalconPortal.Find(&hostgroup); dt.Error != nil {                                               
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }

    hostgroup.Name = inputs.NewName 
    uhostgroup := map[string]interface{}{
        "grp_name":    hostgroup.Name,                                                                   
        "create_user": hostgroup.CreateUser,                                                             
        "come_from":   hostgroup.ComeFrom,
    }   
    dt := db.FalconPortal.Model(&hostgroup).Where("id = ?", hostgroup.ID).Update(uhostgroup)                          
    if dt.Error != nil {                                                                                 
        h.JSONR(c, http.StatusBadRequest, dt.Error)                                                                  
        return
    }   
    h.JSONR(c, "hostgroup profile updated")                                                              
    return
}

func DeleteService(c *gin.Context) {
    grpName := c.Params.ByName("name")
    if grpName == "" {
        h.JSONR(c, http.StatusBadRequest, "grp name is missing")
        return
    }

    hostgroup := fpm.HostGroup{Name: grpName}                                                           
    if dt := db.FalconPortal.Where("grp_name = ?", hostgroup.Name).First(&hostgroup); dt.Error != nil {                                           
        h.JSONR(c, http.StatusBadRequest, dt.Error)                                                              
        return
    }

    tx := db.FalconPortal.Begin()

    //delete hostgroup referance of grp_host table
    if dt := tx.Where("grp_id = ?", hostgroup.ID).Delete(&fpm.GrpHost{}); dt.Error != nil {                       
        h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("delete grp_host got error: %v", dt.Error))                  
        dt.Rollback()                                                                                    
        return
    }
    //delete plugins of hostgroup
    if dt := tx.Where("grp_id = ?", hostgroup.ID).Delete(&fpm.Plugin{}); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("delete plugins got error: %v", dt.Error))
        dt.Rollback()
        return
    }
    //delete aggreators of hostgroup
    if dt := tx.Where("grp_id = ?", hostgroup.ID).Delete(&fpm.Cluster{}); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("delete aggreators got error: %v", dt.Error))
        dt.Rollback()
        return
    }
    //finally delete hostgroup
    if dt := tx.Delete(&fpm.HostGroup{ID: hostgroup.ID}); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        tx.Rollback()
        return
    }

    tx.Commit()

    h.JSONR(c, fmt.Sprintf("hostgroup:%v has been deleted", hostgroup.ID))
    return
}

func ModifyServiceInstance(c *gin.Context) {
    var inputs APIPatchHostGroupHost
    if err := c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }   

    grpName := c.Params.ByName("name")
    if grpName == "" {
        h.JSONR(c, http.StatusBadRequest, "grp name is missing")
        return
    }

    hostgroup := fpm.HostGroup{Name: grpName}
    if dt := db.FalconPortal.Where("grp_name = ?", hostgroup.Name).First(&hostgroup); dt.Error != nil {                                           
        h.JSONR(c, http.StatusBadRequest, dt.Error) 
        return
    }

    action := inputs.Action
    if action != "add" && action != "remove" {
        h.JSONR(c, http.StatusBadRequest, "action must be add or remove")
        return
    }

    switch action {
    case "add":
        bindHostToHostGroup(c, hostgroup, inputs.Hosts)
        return
    case "remove":
        unbindHostToHostGroup(c, hostgroup, inputs.Hosts)
        return
    }
}

func bindHostToHostGroup(c *gin.Context, hostgroup fpm.HostGroup, hosts []string) {                        
    tx := db.FalconPortal.Begin()

    var bindHosts []string                                                                               
    var existHosts []string
    for _, host := range hosts {
        ahost := fpm.Host{Hostname: host}                                                                  
        var id int64
        var ok bool
        if id, ok = ahost.Existing(); !ok {
            if dt := tx.Save(&ahost); dt.Error != nil {
                h.JSONR(c, http.StatusExpectationFailed, dt.Error)
                tx.Rollback()
                return
            }
            id = ahost.ID                                                                   
        }
    
        tGrpHost := fpm.GrpHost{GrpID: hostgroup.ID, HostID: id}
        if ok = tGrpHost.Existing(); ok {
            existHosts = append(existHosts, host)                                          
        } else {
            bindHosts = append(bindHosts, host)
            if dt := tx.Debug().Create(&tGrpHost); dt.Error != nil { 
                h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("create grphost got error: %s , grp_id: %v, host_id: %v", dt.Error, hostgroup.ID, id)) 
                tx.Rollback()
                return
            }                                    
        }
    }

    tx.Commit()

    h.JSONR(c, fmt.Sprintf("%v bind to hostgroup: %s, %v have been exist", bindHosts, hostgroup.Name, existHosts))
    return                                                                                               
}

func unbindHostToHostGroup(c *gin.Context, hostgroup fpm.HostGroup, hosts []string) {                      
    tx := db.FalconPortal.Begin()

    var unbindHosts []string
    for _, host := range hosts {
        dhost := fpm.Host{Hostname: host}
        var id int64
        var ok bool
        if id, ok = dhost.Existing(); ok {
            unbindHosts = append(unbindHosts, host)
        } else {
            log.Debugf("Host %s does not exists!", host)
            continue
        }
        if dt := db.FalconPortal.Where("grp_id = ? AND host_id = ?", hostgroup.ID, id).Delete(&fpm.GrpHost{}); dt.Error != nil {
            h.JSONR(c, http.StatusExpectationFailed, dt.Error)
            tx.Rollback()
            return
        }
    }

    tx.Commit()

    h.JSONR(c, fmt.Sprintf("%v unbind to hostgroup: %s", unbindHosts, hostgroup.Name))
    return
}

func GetAggregatorListOfService(c *gin.Context) {
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

    grpName := c.Params.ByName("name")
    if grpName == "" {
        h.JSONR(c, http.StatusBadRequest, "grp name is missing")
        return
    }   

    hostgroup := fpm.HostGroup{Name: grpName}                                                           
    if dt := db.FalconPortal.Where("grp_name = ?", hostgroup.Name).First(&hostgroup); dt.Error != nil {                                           
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    aggregators := []fpm.Cluster{}
    var dt *gorm.DB
    if limit != -1 && page != -1 {
        dt = db.FalconPortal.Raw(fmt.Sprintf("SELECT * from cluster WHERE grp_id = %d limit %d,%d", hostgroup.ID, page, limit)).Scan(&aggregators)
    } else {
        dt = db.FalconPortal.Where("grp_id = ?", hostgroup.ID).Find(&aggregators)
    }
    if dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }

    hostgroupName := hostgroup.Name
    if len(aggregators) != 0 {
        hostgroupName, err = aggregators[0].HostGroupName()
        if err != nil {
            h.JSONR(c, http.StatusBadRequest, err)
            return
        }
    }

    h.JSONR(c, map[string]interface{}{
        "hostgroup":   hostgroupName,
        "aggregators": aggregators,
    })
    return
}

func GetAggregator(c *gin.Context) {
    aggIDtmp := c.Params.ByName("id")
    if aggIDtmp == "" {
        h.JSONR(c, http.StatusBadRequest, "agg id is missing")
        return
    }
    aggID, err := strconv.Atoi(aggIDtmp)
    if err != nil {
        log.Debugf("aggIDtmp: %v", aggIDtmp)
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }
    aggregator := fpm.Cluster{ID: int64(aggID)}
    if dt := db.FalconPortal.Find(&aggregator); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }
    h.JSONR(c, aggregator)
    return
}

func CreateAggregator(c *gin.Context) {
    var inputs APICreateAggregatorInput                                                                  
    if err := c.Bind(&inputs); err != nil {                                                              
        h.JSONR(c, http.StatusBadRequest, fmt.Sprintf("binding error: %v", err))                                     
        return
    }

    hostgroup := fpm.HostGroup{Name: inputs.GroupName}
    if dt := db.FalconPortal.Where("grp_name = ?", hostgroup.Name).First(&hostgroup); dt.Error != nil {                                           
        h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("find hostgroup error: %v", dt.Error.Error()))           
        return
    }

    agg := fpm.Cluster{                                                                                
        GrpId:       hostgroup.ID,                                                                       
        Numerator:   inputs.Numerator,
        Denominator: inputs.Denominator,
        Endpoint:    inputs.Endpoint,
        Metric:      inputs.Metric,
        Tags:        inputs.Tags,
        DsType:      "GAUGE",
        Step:        inputs.Step,
        Creator:     "root"}
    if dt := db.FalconPortal.Create(&agg); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("create aggregator got error: %v", dt.Error.Error()))
        return
    }
    h.JSONR(c, agg)
    return
}

func UpdateAggregator(c *gin.Context) {
    var inputs APIUpdateAggregatorInput 
    if err := c.Bind(&inputs); err != nil {                                                              
        h.JSONR(c, http.StatusBadRequest, err)                                                                       
        return
    }

    aggregator := fpm.Cluster{ID: inputs.ID}
    if dt := db.FalconPortal.Find(&aggregator); dt.Error != nil {                                              
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)                                                                
        return
    }

    uaggregator := map[string]interface{}{
        "Numerator":   inputs.Numerator,
        "Denominator": inputs.Denominator,
        "Endpoint":    inputs.Endpoint,
        "Metric":      inputs.Metric,
        "Tags":        inputs.Tags,
        "Step":        inputs.Step}
    if dt := db.FalconPortal.Model(&aggregator).Where("id = ?", aggregator.ID).Update(uaggregator).Find(&aggregator); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, dt.Error)
        return
    }
    h.JSONR(c, aggregator)
    return
}

func DeleteAggregator(c *gin.Context) {
    aggIDtmp := c.Params.ByName("id")
    if aggIDtmp == "" {
        h.JSONR(c, http.StatusBadRequest, "agg id is missing")
        return
    }
    aggID, err := strconv.Atoi(aggIDtmp)
    if err != nil {
        log.Debugf("aggIDtmp: %v", aggIDtmp)
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }
    aggregator := fpm.Cluster{ID: int64(aggID)}
    if dt := db.FalconPortal.Find(&aggregator); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("find aggregator got error: %v", dt.Error.Error()))
        return
    }

    if dt := db.FalconPortal.Table("cluster").Where("id = ?", aggID).Delete(&aggregator); dt.Error != nil {
        h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("delete aggregator got error: %v", dt.Error))
        return
    }
    h.JSONR(c, fmt.Sprintf("aggregator:%v has been deleted", aggID))
    return
}

func GetTemplateOfService(c *gin.Context) {
    grpName := c.Params.ByName("name")
    if grpName == "" {
        h.JSONR(c, http.StatusBadRequest, "grp id is missing")
        return
    }

    hostgroup := fpm.HostGroup{Name: grpName}
    if dt := db.FalconPortal.Where("grp_name = ?", hostgroup.Name).First(&hostgroup); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    grpTpls := []fpm.GrpTpl{}
    Tpls := []fpm.Template{}
    db.FalconPortal.Where("grp_id = ?", hostgroup.ID).Find(&grpTpls)
    if len(grpTpls) != 0 { 
        tips := []int64{}
        for _, t := range grpTpls {
            tips = append(tips, t.TplID)
        }   
        tipsStr, _ := u.ArrInt64ToString(tips)
        db.FalconPortal.Where(fmt.Sprintf("id in (%s)", tipsStr)).Find(&Tpls)
    }   
    h.JSONR(c, map[string]interface{}{
        "hostgroup": hostgroup,
        "templates": Tpls,
    })  
    return
}

func BindTemplateToService(c *gin.Context) {
    var inputs APIBindTemplateToGroupInputs
    if err := c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)                                                                       
        return
    }

    hostgroup := fpm.HostGroup{Name: inputs.GroupName}
    if dt := db.FalconPortal.Where("grp_name = ?", hostgroup.Name).First(&hostgroup); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    grpTpl := fpm.GrpTpl{
        GrpID: hostgroup.ID,                                                                             
        TplID: inputs.TplID,                                                                             
    }

    db.FalconPortal.Where("grp_id = ? and tpl_id = ?", hostgroup.ID, inputs.TplID).Find(&grpTpl)               
    // TODO:: how can you do this
    if grpTpl.BindUser != "" {                                                                           
        h.JSONR(c, http.StatusBadRequest, errors.New("this binding already existing, reject!"))                      
        return
    } 

    grpTpl.BindUser = "root"
    if dt := db.FalconPortal.Save(&grpTpl); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }
    h.JSONR(c, grpTpl)
    return
}

func UnbindTemplateToService(c *gin.Context) {
    var inputs APIUnBindTemplateToGroupInputs
    if err := c.Bind(&inputs); err != nil {
        h.JSONR(c, http.StatusBadRequest, err)
        return
    }

    hostgroup := fpm.HostGroup{Name: inputs.GroupName}
    if dt := db.FalconPortal.Where("grp_name = ?", hostgroup.Name).Find(&hostgroup); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    grpTpl := fpm.GrpTpl{
        GrpID: hostgroup.ID,
        TplID: inputs.TplID,
    }

    db.FalconPortal.Where("grp_id = ? and tpl_id = ?", hostgroup.ID, inputs.TplID).Find(&grpTpl)

    if dt := db.FalconPortal.Where("grp_id = ? and tpl_id = ?", hostgroup.ID, inputs.TplID).Delete(&grpTpl); dt.Error != nil {
        h.JSONR(c, http.StatusBadRequest, dt.Error)
        return
    }

    h.JSONR(c, fmt.Sprintf("template: %v is unbind of HostGroup: %v", inputs.TplID, hostgroup.ID))
    return
}

func Routes(r *gin.Engine) {
    db = model.Con()
    //hostgroup
    serviceapi := r.Group("/v1")
    serviceapi.GET("/service", GetServices)
    serviceapi.GET("/service/:name", GetServiceByName)
    serviceapi.POST("/service", CreateService)
    serviceapi.POST("/service/:name/instance", BindInstanceToService)
    serviceapi.PUT("/service/:name/instance", UnbindAInstanceToService)
    serviceapi.PUT("/service", UpdateService)
    serviceapi.DELETE("/service/:name", DeleteService)
    serviceapi.PATCH("/service/:name/instance", ModifyServiceInstance)

    //aggreator
    serviceapi.GET("/service/:name/aggregators", GetAggregatorListOfService)
    serviceapi.GET("/aggregator/:id", GetAggregator)
    serviceapi.POST("/aggregator", CreateAggregator)
    serviceapi.PUT("/aggregator", UpdateAggregator)
    serviceapi.DELETE("/aggregator/:id", DeleteAggregator)

    //template
    serviceapi.GET("/service/:name/template", GetTemplateOfService)
    serviceapi.POST("/service/:name/template", BindTemplateToService)
    serviceapi.PUT("/service/:name/template", UnbindTemplateToService)
}

