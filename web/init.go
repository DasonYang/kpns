package web

import (
    "os"
    "path/filepath"
    "kpns/database"
)

var (
    TemplatePath string
    dbClient    database.DatabaseClient
)

func Init(dbc database.DatabaseClient) error {
    path, _ := os.Getwd()
    TemplatePath = filepath.Join(path, "kpns/templates")
    dbClient = dbc

    return nil
}