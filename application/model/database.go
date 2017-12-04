package model

import (
    "database/sql"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "github.com/spf13/viper"
)

type DBPool struct {
    FalconPortal *gorm.DB
    Graph        *gorm.DB
    Uic          *gorm.DB
    Dashboard    *gorm.DB
    Alarm        *gorm.DB
}

var (
    pool DBPool
)

func Con() DBPool {
    return pool
}

func SetLogLevel() {
    pool.Uic.LogMode(true)
    pool.Graph.LogMode(true)
    pool.FalconPortal.LogMode(true)
    pool.Dashboard.LogMode(true)
    pool.Alarm.LogMode(true)
}

func InitDB() (err error) {
    var fpdb *sql.DB
    portal, err := gorm.Open("mysql", viper.GetString("db.falcon_portal"))
    portal.Dialect().SetDB(fpdb)
    portal.LogMode(true)
    if err != nil {
        return fmt.Errorf("connect to falcon_portal: %s", err.Error())
    }
    portal.SingularTable(true)
    pool.FalconPortal = portal

    var graphdb *sql.DB
    graph, err := gorm.Open("mysql", viper.GetString("db.graph"))
    graph.Dialect().SetDB(graphdb)
    graph.LogMode(true)
    if err != nil {
        return fmt.Errorf("connect to graph: %s", err.Error())
    }
    graph.SingularTable(true)
    pool.Graph = graph

    var uicdb *sql.DB
    uic, err := gorm.Open("mysql", viper.GetString("db.uic"))
    uic.Dialect().SetDB(uicdb)
    uic.LogMode(true)
    if err != nil {
        return fmt.Errorf("connect to uic: %s", err.Error())
    }
    uic.SingularTable(true)
    pool.Uic = uic

    var dashboarddb *sql.DB
    dashboard, err := gorm.Open("mysql", viper.GetString("db.dashboard"))
    dashboard.Dialect().SetDB(dashboarddb)
    dashboard.LogMode(true)
    if err != nil {
        return fmt.Errorf("connect to dashboard: %s", err.Error())
    }
    dashboard.SingularTable(true)
    pool.Dashboard = dashboard

    var alarmdb *sql.DB
    alarm, err := gorm.Open("mysql", viper.GetString("db.alarms"))
    alarm.Dialect().SetDB(alarmdb)
    alarm.LogMode(true)
    if err != nil {
        return fmt.Errorf("connect to alarms: %s", err.Error())
    }
    alarm.SingularTable(true)
    pool.Alarm = alarm

    SetLogLevel()
    return
}

func CloseDB() (err error) {
    err = pool.FalconPortal.Close()
    if err != nil {
        return
    }

    err = pool.Graph.Close()
    if err != nil {
        return
    }

    err = pool.Uic.Close()
    if err != nil {
        return
    }

    err = pool.Dashboard.Close()
    if err != nil {
        return
    }

    err = pool.Alarm.Close()
    if err != nil {
        return
    }

    return
}
