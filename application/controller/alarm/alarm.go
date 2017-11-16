package alarm

import (
    "net/http"

    "github.com/gin-gonic/gin"                                                                           
)

func Routes(r *gin.Engine) {
    alarmapi := r.Group("/v1/alarm")

    alarmapi.GET("/case", GetAlarms)
    alarmapi.GET("/event", GetEvents)
    alarmapi.GET("/note", GetNotesOfAlarm)
}

func GetAlarms(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get alarm list",
    })
}

func GetEvents(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get event list",
    })
}

func GetNotesOfAlarm(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get note list of certain alarm",
    })
}

