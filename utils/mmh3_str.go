package utils

import (
    "log"
    "encoding/binary"

    "github.com/spaolacci/murmur3"
)

func Mmh3py(data string) string {
    var h1, h2 uint64
    // var rh1, rh2 uint64
    // h1, h2 = murmur3.Sum128(buffer.Bytes())
    h1, h2 = murmur3.Sum128([]byte(data))
    log.Printf("h1 = %v, h2 = %v\n", h1, h2)

    // rh1 = divreverse(h1)
    // rh2 = divreverse(h2)

    // log.Printf("rh1 = %v, rh2 = %v\n", rh1, rh2)

    b1 := make([]byte, 8)
    binary.LittleEndian.PutUint64(b1, h1)

    b2 := make([]byte, 8)
    binary.LittleEndian.PutUint64(b2, h2)

    log.Printf("b1 = %v, b2 = %v\n", b1, b2)

    b := append(b1, b2...)

    log.Printf("b = %v\n", string(b))

    var pstartbyte []byte = b
    var pendbyte []byte = b[15:]        /* 16-1 */
    var numsignificantbytes int;        /* number of bytes that matter */
    //var ndigits int;                    /* number of Python long digits */
    //var idigit int = 0;                 /* next free index in v->ob_digit */

    log.Printf("pstartbyte = %v, pendbyte = %v\n", string(pstartbyte), string(pendbyte))

    var is_signed = pendbyte[0] >= 0x80
    log.Printf("is_signed = %v\n", is_signed)
    {
        var p []byte = pendbyte
        var i int
        var pincr int = -1
        var insignficant byte
        if is_signed {
            insignficant = 0xff
        } else{
            insignficant = 0x00
        }

        log.Printf("p = %v, i = %v, pincr = %v, insignficant = %v\n", p, i, pincr, insignficant)

        for i := 0; i < 16; i++  {
            log.Printf("i = %v\n", i)
            if p[i] != insignficant {
                break
            }
        }
        numsignificantbytes = 16 - i
        /* 2's-comp is a bit tricky here, e.g. 0xff00 == -0x0100, so
           actually has 2 significant bytes.  OTOH, 0xff0001 ==
           -0x00ffff, so we wouldn't *need* to bump it there; but we
           do for 0xffff = -0x0001.  To be safe without bothering to
           check every case, bump it regardless. */
        if (is_signed && numsignificantbytes < 16) {
            numsignificantbytes++
        }
    }

    return string(b)
}