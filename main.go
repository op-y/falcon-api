package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "falcon-api/application/controller"
    
    "github.com/gin-gonic/gin"
)

func main() {
    routes := gin.Default()
    go controller.StartGin(":8080", routes)

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigs
        fmt.Println()
        os.Exit(0)
    }() 
    select {}
}
