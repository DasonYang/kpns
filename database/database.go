package database

import (
	"kpns/config"
	"kpns/database/mongo"
)

type KData struct {
	Data map[string]interface{}
}

type ClientData struct {
	Key   string                 `json:"key" binding:"required"`
	Value map[string]interface{} `json:"value"`
}

type DatabaseClient interface {
	Write(string, string, map[string]interface{}, map[string]interface{}) error
	ReadAll(string, string, map[string]interface{}, map[string]interface{}) ([]map[string]interface{}, int)
	ReadOne(string, string, map[string]interface{}) map[string]interface{}
	Delete(string, string, map[string]interface{}) error
	Count(string, string, map[string]interface{}) int
	BulkWrite(string, string, []interface{}) error
	Update(string, string, map[string]interface{}, map[string]interface{}, map[string]interface{}) error
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
