package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func SHA256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}
