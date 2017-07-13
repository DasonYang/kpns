package kpns

import (
    "kpns/config"
)

var (
    Configs config.CfgYaml
    QueueNotification chan PushNotification
)