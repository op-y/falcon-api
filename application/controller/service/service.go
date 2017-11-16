package service

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
    //hostgroup
    serviceapi := r.Group("/v1")
    serviceapi.GET("/service", GetServices)
    serviceapi.GET("/service/:id", GetService)
    serviceapi.POST("/service", CreateService)
    serviceapi.POST("/service/:id/instance", BindInstanceToService)
    serviceapi.PUT("/service/:id/instance", UnbindAInstanceToService)
    serviceapi.PUT("/service", UpdateService)
    serviceapi.DELETE("/service/:id", RemoveService)
    serviceapi.PATCH("/service/:id/instance", ModifyServiceInstance)

    //aggreator
    serviceapi.GET("/service/:id/aggregators", GetAggregatorListOfService)
    serviceapi.GET("/aggregator/:id", GetAggregator)
    serviceapi.POST("/aggregator", CreateAggregator)
    serviceapi.PUT("/aggregator", UpdateAggregator)
    serviceapi.DELETE("/aggregator/:id", DeleteAggregator)

    //template
    serviceapi.GET("/service/:id/template", GetTemplateOfService)
    serviceapi.POST("/service/:id/template", BindTemplateToService)
    serviceapi.PUT("/service/:id/template", UnbindTemplateToService)
}

func GetServices(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get service list",
    })
}

func GetService(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get service",
    })
}

func CreateService(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "create service",
    })
}

func BindInstanceToService(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "bind instances to service",
    })
}

func UnbindAInstanceToService(c *gin.Context) {
    c.JSON(http.StatusOK,gin.H{
        "message": "unbind instances to service",
    })
}

func UpdateService(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "update service",
    })
}

func RemoveService(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "remove service",
    })
}

func ModifyServiceInstance(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "modify service",
    })
}

func GetAggregatorListOfService(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get aggregator list of service",
    })
}

func GetAggregator(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get aggregator",
    })
}

func CreateAggregator(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "create aggregator",
    })
}

func UpdateAggregator(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "update aggregator",
    })
}

func DeleteAggregator(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "delete aggregator",
    })
}

func GetTemplateOfService(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get template of service",
    })
}

func BindTemplateToService(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "bind template to service",
    })
}

func UnbindTemplateToService(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "unbind template to service",
    })
}

