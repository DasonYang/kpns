package database

import (
    "kpns/config"
    "kpns/database/mongo"
)

type KData struct {
    Data    map[string]interface{}
}

type ClientData struct {
    Key     string                      `json:"key" binding:"required"`
    Value   map[string]interface{}      `json:"value"`
}

type DatabaseClient interface {
    Write(string, string, map[string]interface{}) error
    Read(string, string)  map[string]interface{}
}

func NewDB(cfg config.CfgYaml) DatabaseClient {
    var db DatabaseClient = nil
    switch cfg.Core.Database {
    case "mongo":
        db = mongo.New()
    default:
        db = mongo.New()
    }
    return db
}