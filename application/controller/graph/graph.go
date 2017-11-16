package graph

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
    graphapi := r.Group("/v1/graph")

    graphapi.GET("/endpoint-object", GetEndpointObject)
    graphapi.GET("/endpoint", GetEndpointByRegExp)
    graphapi.GET("/endpoint-counter", GetEndpointCounterByRegExp)
    graphapi.POST("/history", GetGraphDrawData)
    graphapi.POST("/last-point", GetGraphLastPoint)
    graphapi.DELETE("/endpoint", DeleteGraphEndpoint)
    graphapi.DELETE("/counter", DeleteGraphCounter)
}

func GetEndpointObject(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get endpoint object",
    })
}

func GetEndpointByRegExp(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get endpoint by regular expression",
    })
}

func GetEndpointCounterByRegExp(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get endpoint counter by regular expression",
    })
}

func GetGraphDrawData(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get graph draw data",
    })
}

func GetGraphLastPoint(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get graph last point",
    })
}

func DeleteGraphEndpoint(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "delete graph endpoint",
    })
}

func DeleteGraphCounter(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "delete graph counter",
    })
}

