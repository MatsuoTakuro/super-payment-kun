package pkg

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (rc RealClocker) Now() time.Time {
	return time.Now()
}

// FixedClocker is a clocker that always returns the same time.
// It is useful for testing.
type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2024, 5, 26, 12, 0, 0, 0, time.UTC)
}
