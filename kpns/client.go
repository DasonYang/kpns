package kpns

import (
    
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
