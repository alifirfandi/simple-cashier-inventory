package helper

import (
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandString
//
// Ref: How to generate a random string of a fixed length in Go?
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// SplitLastStr
//
// Ref: Split a string at the last occurrence of the separator in golang
// https://stackoverflow.com/questions/51598250/split-a-string-at-the-last-occurrence-of-the-separator-in-golang
func SplitLastStr(s, sep string) (firstSub string, lastSub string) {
	lastInd := strings.LastIndex(s, sep)
	if lastInd < 0 {
		firstSub = ""
		lastSub = ""
	}
	firstSub = s[:lastInd]
	lastSub = s[lastInd+1:]
	return
}
