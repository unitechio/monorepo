package date

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

func Format(t time.Time, full bool) string {
	if full {
		return t.Format("02/01/2006 15:04:05")
	}
	return t.Format("02/01/2006")
}

func AddDays(t time.Time, d int) time.Time {
	return t.AddDate(0, 0, d)
}

func AddMonths(t time.Time, m int) time.Time {
	return t.AddDate(0, m, 0)
}

func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	return StartOfDay(t).Add(24*time.Hour - time.Nanosecond)
}

func Now() time.Time {
	return time.Now()
}

func FormatDate(t time.Time, full bool) string {
	if full {
		return t.Format("02/01/2006 15:04:05")
	}
	return t.Format("02/01/2006")
}

func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func SHA256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}
