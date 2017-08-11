package web

import (
    "os"
    "fmt"
    "path/filepath"
    "kpns/database"
)

var (
    TemplatePath string
    dbClient    database.DatabaseClient
)

func Init(dbc database.DatabaseClient) error {
    path, _ := os.Getwd()
    TemplatePath = filepath.Join(path, "web/templates")
    fmt.Printf("TemplatePath = %v\n", TemplatePath)
    dbClient = dbc

    return nil
}