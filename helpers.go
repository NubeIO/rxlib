package rxlib

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateCustomUUID(prefix ...string) string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	hexString := hex.EncodeToString(b)
	if len(prefix) > 0 {
		return fmt.Sprintf("%s_%s", prefix[0], hexString)
	}
	return hexString
}

func TimeSince(t time.Time) string {
	var duration = time.Since(t)
	switch {
	case duration < 30*time.Second:
		return "just now"
	case duration < 1*time.Minute:
		return fmt.Sprintf("%d sec", int(duration.Seconds()))
	case duration < 1*time.Hour:
		return fmt.Sprintf("%d min ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < 30*24*time.Hour: // Approximating a month as 30 days
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	case duration < 365*24*time.Hour: // Approximating a year as 365 days
		return fmt.Sprintf("%d months ago", int(duration.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%d years ago", int(duration.Hours()/(24*365)))
	}
}
