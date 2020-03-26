package util

import (
	"math/rand"
	"time"
)

func RandomName(n int) string{
	var letters  = []byte("abcdefghijklmnopqrstuvwxyz1234567890")
	result:=make([]byte,n)
	rand.Seed(time.Now().Unix())
	for i:=range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return  string(result)
}
