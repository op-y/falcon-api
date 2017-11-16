package template

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
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

func GetTemplates(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get template list",
    })
}

func GetTemplate(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get template",
    })
}

func CreateTemplate(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "create template",
    })
}

func UpdateTemplate(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "update template",
    })
}

func DeleteTemplate(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "delete template",
    })
}

func CreateActionToTmplate(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "create action to template",
    })
}

func UpdateActionToTmplate(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "update action to template",
    })
}

func GetActionByID(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get action by ID",
    })
}

