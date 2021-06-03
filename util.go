package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func ToBase62(number int) string {

	divisor := 62
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var shortURL []rune

	for number > 0 {

		shortURL = append(shortURL, chars[number%divisor])
		number = number / divisor
	}

	return reverse(string(shortURL))

}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func computeSHA256Base64(input string) string {

	h := sha256.New()
	h.Write([]byte(input))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	s := base64.StdEncoding.EncodeToString([]byte(hash))
	return s
}

func substring(input string) string {
	chars := []int32(input)

	return string(chars[:6])
}

func GenerateShortURLCodeOld(input string) string {

	shortURLCode := substring(computeSHA256Base64(input))
	return shortURLCode
}
