package database

import (
)

type KData struct {
    Data    map[string]interface{}
}

type ClientData struct {
    Key     string                      `json:"key" binding:"required"`
    Value   map[string]interface{}      `json:"value"`
}

type DatabaseWrapper interface {
    Init() error
    Write() error
    Read()
    Update()
    Delete() error
}