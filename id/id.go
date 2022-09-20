package id

import (
	"crypto/rand"
	"encoding/base64"
)

func withIDSuffix(input string) string {
	randBytes := make([]byte, 8)
	_, _ = rand.Read(randBytes)
	randStr := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes)

	return input + randStr
}

func NewRunID() string {
	return withIDSuffix("run_")
}

func NewTestID() string {
	return withIDSuffix("test_")
}
