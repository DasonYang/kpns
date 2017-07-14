package kpns

import (
    "kpns/config"

    apns "github.com/sideshow/apns2"
)

var (
    Configs config.CfgYaml
    QueueNotification chan PushNotification
    ApnsClient *apns.Client
)