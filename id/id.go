package id

import (
	"encoding/base64"
	"math/rand"
	"time"
)

func withIdSuffix(input string) string {
	rand.Seed(time.Now().Unix())

	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	randStr := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes)

	return input + randStr
}

func NewRunId() string {
	return withIdSuffix("run_")
}

func NewTestId() string {
	return withIdSuffix("test_")
}
