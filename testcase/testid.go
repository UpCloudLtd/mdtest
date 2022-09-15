package testcase

import (
	"encoding/base64"
	"math/rand"
	"time"
)

func testId() string {
	rand.Seed(time.Now().Unix())

	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	randStr := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes)

	return randStr
}
