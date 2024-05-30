package pkg

import "time"

const (
	DateLayout = "2006-01-02"
)

// TruncateToDate truncates the time part, leaving only the date part.
func TruncateToDate(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}

func FormatToDate(t time.Time) string {
	return t.Format(DateLayout)
}
