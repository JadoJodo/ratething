package thing

import (
	"github.com/google/uuid"
)

type Thing struct {
	ID       int64
	UUID     uuid.UUID
	Name     string
	Category string
}
