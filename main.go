package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "falcon-api/application/controller"
    gclient "falcon-api/application/graph"
    "falcon-api/application/model"
    
    log "github.com/Sirupsen/logrus"
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
)

var cluster map[string]string


func main() {
    viper.AddConfigPath(".")
    viper.SetConfigName("config")

    err := viper.ReadInConfig()
    if err != nil {
        log.Fatal(err.Error())
    }

    level := viper.GetString("log_level")
    switch level {
    case "info":
        log.SetLevel(log.InfoLevel)
    case "debug":
        log.SetLevel(log.DebugLevel)
    case "warn":
        log.SetLevel(log.WarnLevel)
    default:
        log.Fatal("log conf only allow [info, debug, warn], please check your confguire")
    }   

    gclient.Start(viper.GetStringMapString("graphs.cluster"))

    if err := model.InitDB(); err != nil {
        log.Fatalf("db conn failed with error %s", err.Error())
    } 


    if viper.GetString("log_level") != "debug" {
        gin.SetMode(gin.ReleaseMode)
    }
    //gin.DisableConsoleColor()
    //f, _ := os.Create("gin.log")
    //gin.DefaultWriter = io.MultiWriter(f)
    //gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

    routes := gin.Default()
    log.Debugf("will start with port:%v", viper.GetString("web_port"))
    go controller.StartGin(viper.GetString("web_port"), routes)

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigs
        fmt.Println()
        os.Exit(0)
    }() 
    select {}
}
