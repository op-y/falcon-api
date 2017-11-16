package strategy

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
    strategyapi := r.Group("/v1/strategy")

    strategyapi.GET("", GetStrategys)
    strategyapi.GET("/:id", GetStrategy)
    strategyapi.POST("", CreateStrategy)
    strategyapi.PUT("", UpdateStrategy)
    strategyapi.DELETE("/:id", DeleteStrategy)

    metricapi := r.Group("/v1/metric")
    metricapi.GET("", QueryMetrics)
}

func GetStrategys(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get strategy list",
    })
}

func GetStrategy(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get strategy",
    })
}

func CreateStrategy(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "create strategy",
    })
}

func UpdateStrategy(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "update strategy",
    })
}

func DeleteStrategy(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "delete strategy",
    })
}

func QueryMetrics(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "query metrics",
    })
}

