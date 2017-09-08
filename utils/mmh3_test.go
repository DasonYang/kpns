package utils

import (
    "fmt"
    "testing"
    
    "github.com/reusee/mmh3"
)

func TestOfficialMmh3(t *testing.T) {
    var str = "ANd850e608518f6b51094dcom.easyn.EasyN_P1android"

    fmt.Printf("Official mmh3 = %v\n", mmh3.Sum128([]byte(str)))

    fmt.Printf("my mmh3 = %v\n", Hash128(str))
}