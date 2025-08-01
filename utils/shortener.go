package utils

import (
	"crypto/sha256"
	"math/big"
	"strconv"
	"sync/atomic"
	"time"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	counter int64  = time.Now().Unix()
	salt    string = "secret-salt"
)

func toBase62(num *big.Int) string {
	if num.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}

	var result []byte
	base := big.NewInt(62)
	zero := big.NewInt(0)

	n := new(big.Int).Set(num)

	for n.Cmp(zero) > 0 {
		remainder := new(big.Int)
		n.DivMod(n, base, remainder)
		result = append([]byte{base62Chars[remainder.Int64()]}, result...)
	}

	return string(result)
}

func GenerateShortUrl(url string) string {

	currentCounter := atomic.AddInt64(&counter, 1)
	toHash := url + salt + strconv.FormatInt(currentCounter, 10)

	algorithm := sha256.New()
	algorithm.Write([]byte(toHash))
	urlHashBytes := algorithm.Sum(nil)

	hashSubset := urlHashBytes[:8]
	generatedNumber := new(big.Int).SetBytes(hashSubset)

	finalString := toBase62(generatedNumber)

	if len(finalString) >= 6 {
		return finalString[:6]
	}
	return finalString
}
