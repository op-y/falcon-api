package dashboard

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
    dashboardapi := r.Group("/v1/dashboard")

    dashboardapi.GET("/tmpgraph/:id", GetTmpGraph)
    dashboardapi.POST("/tmpgraph", CreateTmpGraph)

    dashboardapi.GET("/graph/:id", GetDashboardGraph)
    dashboardapi.POST("/graph", CreateDashboardGraph)
    dashboardapi.PUT("/graph/:id", UpdateDashboardGraph)
    dashboardapi.DELETE("/graph/:id", DeleteDashboardGraph)

    dashboardapi.GET("/graphs/screen/:id", GetDashboardGraphsByScreenID)
}

func GetTmpGraph(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get temporary graph",
    })
}

func CreateTmpGraph(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "create temporary graph",
    })
}

func GetDashboardGraph(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get dashboard graph",
    })
}

func CreateDashboardGraph(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "create dashboard graph",
    })
}

func UpdateDashboardGraph(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "update dashboard graph",
    })
}

func DeleteDashboardGraph(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "delete dashboard graph",
    })
}

func GetDashboardGraphsByScreenID(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get dashboard graph by screen ID",
    })
}

