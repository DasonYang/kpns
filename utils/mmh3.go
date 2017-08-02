package utils

import (
    "log"
    "bytes"
    "encoding/binary"

    "github.com/spaolacci/murmur3"
)

const (
    LongShift = 30
    LongBase = uint32(1) << LongShift
    LongMash = uint32(LongBase - 1)
    LongDecimalShift = 9
    LongDecimalBase = uint32(1000000000)
)

func Sum128toString(udid string, appid string, platform string) string {
    var buffer bytes.Buffer
    buffer.WriteString(udid)
    buffer.WriteString(appid)
    buffer.WriteString(platform)
    var key string = "##"
    var v []uint64

    h1, h2 := murmur3.Sum128(buffer.Bytes())

    b1 := make([]byte, 8)
    binary.LittleEndian.PutUint64(b1, h1)

    b2 := make([]byte, 8)
    binary.LittleEndian.PutUint64(b2, h2)

    b := append(b1, b2...)

    var pstartbyte []byte = b
    var pendbyte []byte = b[15:]
    var numsignificantbytes int
    var ndigits int
    var idigit int = 0
    var is_signed = pendbyte[0] >= 0x80

    {
        var p []byte = pendbyte
        var i int
        var insignficant byte
        if is_signed {
            insignficant = 0xff
        } else{
            insignficant = 0x00
        }

        for i := 0; i < 16; i++  {
            if p[i] != insignficant {
                break
            }
        }
        numsignificantbytes = 16 - i

        if (is_signed && numsignificantbytes < 16) {
            numsignificantbytes++
        }
    }
    ndigits = (numsignificantbytes * 8 + 30 - 1) / 30
    v = make([]uint64, ndigits)

    {
        var i int;
        var carry uint64 = 1
        var accum uint64 = 0
        var accumbits uint = 0
        var p []byte = pstartbyte

        for i = 0; i < numsignificantbytes; i++ {
            var thisbyte uint64 = uint64(p[i])

            if (is_signed) {
                thisbyte = (0xff ^ thisbyte) + carry
                carry = thisbyte >> 8
                thisbyte &= 0xff
            }

            accum |= uint64(thisbyte << accumbits)

            accumbits += 8;
            if (accumbits >= 30) {
                v[idigit] = uint64(uint32(accum) & LongMash)
                idigit++
                accum >>= 30
                accumbits -= 30
            }
        }
        if accumbits != 0 {
            v[idigit] = uint64(accum)
            idigit++
        }
    }

    {
        var size, strlen, size_a, i, j int
        var rem, tenpow uint32
        var pout []uint64
        var scratch []uint64 = make([]uint64, ndigits)
        var p []byte
        var negative int = 0
        var addL int = 1

        pout = scratch

        size = 0
        size_a = ndigits
        
        for i = size_a-1; i >= 0; i-- {

            var hi uint32 = uint32(v[i])

            for j := 0; j < size; j++ {
                var z uint64 = uint64(pout[j]) << LongShift | uint64(hi)
                hi = uint32(z / uint64(LongDecimalBase))
                pout[j] = uint64(uint32(z) - (hi * LongDecimalBase))
            }

            for hi > 0 {
                pout[size] = uint64(hi % LongDecimalBase)
                size++
                hi /= LongDecimalBase
            }
        }
        
        if size == 0 {
            pout[size] = 0
            size++
        }

        strlen = addL + negative + 1 + (size - 1) * LongDecimalShift
        tenpow = 10
        rem = uint32(pout[size-1])
        for rem >= tenpow {
            tenpow *= 10
            strlen++
        }

        p = make([]byte, strlen-1)
        pos := strlen-2
        // p[pos] = 0x00
        // pos--

        for i=0; i < size - 1; i++ {
            rem = uint32(pout[i])
            for j = 0; j < LongDecimalShift; j++ {
                p[pos] = 0x30 + byte(rem % 10)
                pos--
                rem /= 10
            }
        }

        rem = uint32(pout[i])
        for rem != 0 {
            p[pos] = 0x30 + byte(rem % 10)
            pos--
            rem /= 10
        }

        key = string(p)
    }
    return key
}

func Mmh3py(data string) string {

    h1, h2 := murmur3.Sum128([]byte(data))
    log.Printf("h1 = %v, h2 = %v\n", h1, h2)
    // var incr int = 1
    var v []uint64


    b1 := make([]byte, 8)
    binary.LittleEndian.PutUint64(b1, h1)

    b2 := make([]byte, 8)
    binary.LittleEndian.PutUint64(b2, h2)

    b := append(b1, b2...)

    // log.Printf("b = %v\n", string(b))

    var pstartbyte []byte = b
    var pendbyte []byte = b[15:]        /* 16-1 */
    var numsignificantbytes int;        /* number of bytes that matter */
    var ndigits int;                    /* number of Python long digits */
    var idigit int = 0;                 /* next free index in v->ob_digit */

    // log.Printf("pstartbyte = %v, pendbyte = %v\n", string(pstartbyte), string(pendbyte))

    var is_signed = pendbyte[0] >= 0x80
    // log.Printf("is_signed = %v\n", is_signed)
    {
        var p []byte = pendbyte
        var i int
        // var pincr int = -1
        var insignficant byte
        if is_signed {
            insignficant = 0xff
        } else{
            insignficant = 0x00
        }

        // log.Printf("p = %v, i = %v, pincr = %v, insignficant = %v\n", p, i, pincr, insignficant)

        for i := 0; i < 16; i++  {
            // log.Printf("i = %v\n", i)
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
    ndigits = (numsignificantbytes * 8 + 30 - 1) / 30
    v = make([]uint64, ndigits)
    // log.Printf("ndigits = %v, numsignificantbytes = %v\n", ndigits, numsignificantbytes)

    {
        var i int;
        var carry uint64 = 1                    /* for 2's-comp calculation */
        var accum uint64 = 0                    /* sliding register */
        var accumbits uint = 0             /* number of bits in accum */
        var p []byte = pstartbyte
        // log.Printf("pstartbyte = %v\n", pstartbyte)
        for i = 0; i < numsignificantbytes; i++ {
            var thisbyte uint64 = uint64(p[i])
            // printf("%" PRIu64 "\n", thisbyte);
            /* Compute correction for 2's comp, if needed. */
            if (is_signed) {
                thisbyte = (0xff ^ thisbyte) + carry
                carry = thisbyte >> 8
                thisbyte &= 0xff
            }
             // Because we're going LSB to MSB, thisbyte is
             //   more significant than what's already in accum,
             //   so needs to be prepended to accum. 
            accum |= uint64(thisbyte << accumbits)
            // printf("accum = %d\n", accum);
            accumbits += 8;
            if (accumbits >= 30) {
                /* There's enough to fill a Python digit. */
                v[idigit] = uint64(uint32(accum) & LongMash)
                idigit++
                accum >>= 30
                accumbits -= 30
            }
        }
        if accumbits != 0 {
            v[idigit] = uint64(accum)
            idigit++
        }
    }

    // for idx := 0; idx < idigit; idx++ {
    //     log.Printf("v[%v] = %v\n", idx, v[idx])
    // }

    {
        var size, strlen, size_a, i, j int
        var rem, tenpow uint32
        // var pin []uint64
        var pout []uint64
        var scratch []uint64 = make([]uint64, ndigits)
        var p []byte
        var negative int = 0
        var addL int = 1

        // pin = v
        pout = scratch
        // for idx := 0; idx < idigit; idx++ {
        //     log.Printf("v[%v] = %v\n", idx, pin[idx])
        // }
        size = 0
        size_a = ndigits
        
        for i = size_a-1; i >= 0; i-- {

            var hi uint32 = uint32(v[i])
            // log.Printf("Line = 147, i = %v, hi = %v\n", i, hi)
            for j := 0; j < size; j++ {
                var z uint64 = uint64(pout[j]) << LongShift | uint64(hi)
                // log.Printf("\nLine = 150, z = %v\n", z)
                hi = uint32(z / uint64(LongDecimalBase))
                // log.Printf("Line = 152, hi = %v\n", hi)
                pout[j] = uint64(uint32(z) - (hi * LongDecimalBase))
                // log.Printf("Line = 154, pout[%v] = %v\n", j,pout[j])
            }
            // printf("Line = 233, i = %d, result = %d\n", i, --i >= 0);
            // sleep(1);
            for hi > 0 {
                pout[size] = uint64(hi % LongDecimalBase)
                size++
                hi /= LongDecimalBase
            }
        }
        
        /* pout should have at least one digit, so that the case when a = 0
           works correctly */
        if size == 0 {
            pout[size] = 0
            size++
        }
        // log.Printf("Line = 237, pout = %v\n", pout)
        /* calculate exact length of output string, and allocate */
        strlen = addL + negative + 1 + (size - 1) * LongDecimalShift
        tenpow = 10
        rem = uint32(pout[size-1])
        for rem >= tenpow {
            tenpow *= 10
            strlen++
        }

        /* fill the string right-to-left */
        p = make([]byte, strlen-1)
        pos := strlen-2
        // p = make([]byte, strlen+1)
        // pos := strlen
        // p[pos] = 0x00
        // pos--
        // if (addL > 0) {
        //     p[pos] = 0x4C
        //     pos--
        // }
        // log.Printf("Line = 266, p = %v, pos = %v\n", string(p), pos)
        /* pout[0] through pout[size-2] contribute exactly
           LongDecimalShift digits each */
        for i=0; i < size - 1; i++ {
            rem = uint32(pout[i])
            for j = 0; j < LongDecimalShift; j++ {
                p[pos] = 0x30 + byte(rem % 10)
                pos--
                rem /= 10
            }
        }
        /* pout[size-1]: always produce at least one decimal digit */
        rem = uint32(pout[i])
        for rem != 0 {
            p[pos] = 0x30 + byte(rem % 10)
            pos--
            rem /= 10
        }

        // log.Printf("p = %v\n", string(p))
        return string(p)
    }

    return string(-1)
}