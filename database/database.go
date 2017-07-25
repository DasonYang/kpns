package database

import (
    // "kpns/config"
    // "kpns/database/mongo"
)

type KData struct {
    Data    map[string]interface{}
}

type ClientData struct {
    Key     string                      `json:"key" binding:"required"`
    Value   map[string]interface{}      `json:"value"`
}

type DatabaseClient interface {
    // Init() error
    Write(string, string, map[string]interface{}) error
    Read(string, string)  map[string]interface{}
    // SetDatabase(string) error
    // SetTable(string) error
    // Update()
    // Delete() error
}

// func NewDB(cfg config.CfgYaml) *DatabaseClient {
//     switch cfg.Core.Database {
//     case "mongo":
//     default:
//         var db DatabaseClient = mongo.New()
//     }
//     return nil
// }