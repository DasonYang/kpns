package kpns

import (
    "kpns/config"
    "kpns/database"

    apns "github.com/sideshow/apns2"
)

var (
    Configs config.CfgYaml
    QueueNotification chan PushNotification
    ApnsClient *apns.Client
    DBClient    database.DatabaseClient
)