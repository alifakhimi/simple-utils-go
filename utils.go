package simutils

import (
	"strings"
	"time"
)

func FileNameFriendlyNowTime() string {
	return FileNameFriendlyTime(time.Now().Local())
}

func FileNameFriendlyTime(t time.Time) string {
	ts := t.Format(time.RFC3339)
	return strings.Replace(strings.Replace(ts, ":", "", -1), "-", "", -1)
}
