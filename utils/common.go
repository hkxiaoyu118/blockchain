package utils

import (
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

func StrGetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func StrGetRandStringEx(length int, str string) string {
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func StrSubString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}
	return string(rs[start:end])
}

func StrGetRandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetUnixTimeMs() int64 {
	return time.Now().UnixNano() / 1e6
}

func RemovePrefix0x(s string) string {
	return strings.TrimPrefix(s, "0x")
}

func AddPrefix0x(s string) string {
	if !strings.HasPrefix(s, "0x") {
		return "0x" + s
	}
	return s
}

func SplitByLength(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	length := 0
	for _, r := range s {
		chunk[length] = r
		length++
		if length == chunkSize {
			chunks = append(chunks, string(chunk))
			length = 0
		}
	}
	if length != 0 {
		chunks = append(chunks, string(chunk[:length]))
	}
	return chunks
}

func Base64ToHex(s string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(s)
	if err == nil {
		hexData := hex.EncodeToString(decodedData)
		return hexData, nil
	}
	return "", err
}

func HexToBase64(s string) (string, error) {
	decodedData, err := hex.DecodeString(s)
	if err == nil {
		base64Data := base64.StdEncoding.EncodeToString(decodedData)
		return base64Data, nil
	}
	return "", err
}
