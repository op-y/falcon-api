package controller

import (
    "net/http"

    "falcon-api/application/controller/alarm"
    "falcon-api/application/controller/graph"
    "falcon-api/application/controller/service"
    "falcon-api/application/controller/strategy"
    "falcon-api/application/controller/template"
    "falcon-api/application/controller/uic"

    "github.com/gin-gonic/gin"
)

func StartGin(port string, r *gin.Engine) {
    SystemRoutes(r)
    alarm.Routes(r)
    graph.Routes(r)
    service.Routes(r)
    strategy.Routes(r)
    template.Routes(r)
    uic.Routes(r)
    r.Run(port)
}

func SystemRoutes(r *gin.Engine) {
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, I'm SRE!")
    }) 

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })
}
