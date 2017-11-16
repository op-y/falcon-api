package uic

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
    userapi := r.Group("/v1/user")
    userapi.GET("", GetUsers)
    userapi.GET("/:id", GetUser)
    userapi.POST("", CreateUser)
    userapi.PUT("/:id", UpdateUser)
    userapi.GET("/:id/team", GetUserTeams)
    userapi.DELETE("", DeleteUser)

    teamapi := r.Group("/v1/team")
    teamapi.GET("", GetTeams)
    teamapi.GET("/:id", GetTeamByID)
    teamapi.POST("", CreateTeam)
    teamapi.PUT("", UpdateTeam)
    teamapi.DELETE("/:id", DeleteTeam)
}

func GetUsers(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get user list",
    })
}

func GetUser(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get user",
    })
}

func CreateUser(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "create user",
    })
}

func UpdateUser(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "update user",
    })
}

func GetUserTeams(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get teams of user",
    })
}

func DeleteUser(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "delete user",
    })
}

func GetTeams(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get team list",
    })
}

func GetTeamByID(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "get team by ID",
    })
}

func CreateTeam(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "creat team",
    })
}

func UpdateTeam(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "update team",
    })
}

func DeleteTeam(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "delete team",
    })
}

