package kpns

import (
    // "fmt"
    "log"
    "bytes"
    
    // "encoding/binary"

    "kpns/utils"

    // "github.com/spaolacci/murmur3"
    // "github.com/dmcgowan/mmh3"

    // "github.com/reusee/mmh3"
    // "strconv"
)

func divreverse(n uint64) uint64 {

    var ret uint64
    for i := n; i > 0; {
        ret = ret<<1 + i%2
        i = i / 2
    }

    return ret
}

type ClientData struct {
    Key     string                  `bson:"key"`
    Value   map[string]interface{}  `bson:"value"`
}

func GetClientKey(udid string, appid string, platform string) string {
    var buffer bytes.Buffer
    buffer.WriteString(udid)
    buffer.WriteString(appid)
    buffer.WriteString(platform)
    var key string

    log.Printf("buffer = %v\n", buffer.String())

    // key = utils.Mmh3py(buffer.String())
    key = utils.Mmh3py("abc")

    // rb := make([]byte, 16)

    // for idx, item := range b {
    //     rb[16-idx-1] = item
    //     log.Printf("item = %v\n", item)
    // }

    // log.Printf("rb = %v\n", rb)


    // i := binary.LittleEndian.Uint64(rb)

    // log.Printf("i = %v\n", i)

    log.Printf("key = %v\n", key)
    // ---------github.com/reusee/mmh3------------
    // ret := mmh3.Sum128(buffer.Bytes())
    // log.Printf("ret = %v\n", ret)
    // i, j := strconv.Atoi(string(ret))
    // log.Printf("i = %v, %v\n", i, j)
    // key = fmt.Sprintf("%v", ret)

    return key
}