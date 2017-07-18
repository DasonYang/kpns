package kpns

import(
    "fmt"
)

type PushRequest struct {
    Notifications   []PushNotification  `json:"notifications" binding:"required"`
}

// PushNotification is single notification request
type PushNotification struct {
    // Common
    // Tokens           []string `json:"tokens" binding:"required"`
    // Platform         int      `json:"platform" binding:"required"`
    // Message          string   `json:"message,omitempty"`
    // Title            string   `json:"title,omitempty"`
    // Priority         string   `json:"priority,omitempty"`
    // ContentAvailable bool     `json:"content_available,omitempty"`
    // Sound            string   `json:"sound,omitempty"`
    // Data             D        `json:"data,omitempty"`
    // Retry            int      `json:"retry,omitempty"`
    // wg               *sync.WaitGroup
    // log              *[]LogPushEntry

    // Android
    // APIKey                string           `json:"api_key,omitempty"`
    // To                    string           `json:"to,omitempty"`
    // CollapseKey           string           `json:"collapse_key,omitempty"`
    // DelayWhileIdle        bool             `json:"delay_while_idle,omitempty"`
    // TimeToLive            *uint            `json:"time_to_live,omitempty"`
    // RestrictedPackageName string           `json:"restricted_package_name,omitempty"`
    // DryRun                bool             `json:"dry_run,omitempty"`
    // Notification          fcm.Notification `json:"notification,omitempty"`

    // iOS
    // Expiration     int64    `json:"expiration,omitempty"`
    // ApnsID         string   `json:"apns_id,omitempty"`
    // Topic          string   `json:"topic,omitempty"`
    // Badge          *int     `json:"badge,omitempty"`
    // Category       string   `json:"category,omitempty"`
    // URLArgs        []string `json:"url-args,omitempty"`
    // Alert          Alert    `json:"alert,omitempty"`
    // MutableContent bool     `json:"mutable-content,omitempty"`

    Tokens              []string    `json:"tokens" binding:"required"`
    Platform            int         `json:"platform" binding:"required"`
}

// InitWorkers for initialize all workers.
func InitWorkers(workerNum int64, queueNum int64) {
    // LogAccess.Debug("worker number is ", workerNum, ", queue number is ", queueNum)
    fmt.Printf("worker number is ", workerNum, ", queue number is ", queueNum)
    QueueNotification = make(chan PushNotification, queueNum)
    for i := int64(0); i < workerNum; i++ {
        go startWorker()
    }
}

func startWorker() {
    for {
        notification := <-QueueNotification
        switch notification.Platform {
        case PlatFormIos:
            // PushToIOS(notification)
            fmt.Printf("Push to iOS")
        case PlatFormAndroid:
            // PushToAndroid(notification)
            fmt.Printf("Push to Android")
        }
    }
}

func queueNotification(req PushRequest) int {
    var count int
    newNotification := []PushNotification{}

    for _, notification := range req.Notifications {
        switch notification.Platform {
        case PlatFormIos:
            fmt.Printf("Platfor iOS")
        case PlatFormAndroid:
            fmt.Printf("Platfor Android")
        }
        newNotification = append(newNotification, notification)
    }

    return count
}

// InitAPNSClient use for initialize APNs Client.
// func InitAPNSClient() error {

//     var err error
//     ext := filepath.Ext(PushConf.Ios.KeyPath)

//     switch ext {
//     case ".p12":
//         CertificatePemIos, err = certificate.FromP12File(PushConf.Ios.KeyPath, PushConf.Ios.Password)
//     case ".pem":
//         CertificatePemIos, err = certificate.FromPemFile(PushConf.Ios.KeyPath, PushConf.Ios.Password)
//     default:
//         err = errors.New("wrong certificate key extension")
//     }

//     if err != nil {
//         LogError.Error("Cert Error:", err.Error())

//         return err
//     }

//     if PushConf.Ios.Production {
//         ApnsClient = apns.NewClient(CertificatePemIos).Production()
//     } else {
//         ApnsClient = apns.NewClient(CertificatePemIos).Development()
//     }

//     return nil
// }

// queueNotification add notification to queue list.
// func queueNotification(req RequestPush) (int, []LogPushEntry) {
//     var count int
//     wg := sync.WaitGroup{}
//     newNotification := []PushNotification{}
//     for _, notification := range req.Notifications {
//         switch notification.Platform {
//         case PlatFormIos:
//             if !PushConf.Ios.Enabled {
//                 continue
//             }
//         case PlatFormAndroid:
//             if !PushConf.Android.Enabled {
//                 continue
//             }
//         }
//         newNotification = append(newNotification, notification)
//     }

//     log := make([]LogPushEntry, 0, count)
//     for _, notification := range newNotification {
//         if PushConf.Core.Sync {
//             notification.wg = &wg
//             notification.log = &log
//             notification.AddWaitCount()
//         }
//         QueueNotification <- notification
//         count += len(notification.Tokens)
//     }

//     if PushConf.Core.Sync {
//         wg.Wait()
//     }

//     StatStorage.AddTotalCount(int64(count))

//     return count, log
// }