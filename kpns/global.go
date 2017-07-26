package kpns

import (
    "kpns/config"
    "kpns/database"

    // apns "github.com/sideshow/apns2"

)

var (
    Configs config.CfgYaml
    QueueNotification chan PushNotification
    DBClient    database.DatabaseClient
)