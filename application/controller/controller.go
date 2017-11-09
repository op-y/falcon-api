package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func StartGin(port string, r *gin.Engine) {
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, I'm Hualala SRE!")
    })  
    r.Run(port)
}

