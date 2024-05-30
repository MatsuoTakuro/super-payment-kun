package pkg

import "github.com/google/uuid"

// NewUUID generates a new UUID (v4).
func NewUUID() string {
	return uuid.New().String()
}

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
