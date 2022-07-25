package thing

import (
	"github.com/google/uuid"
)

type Thing struct {
	ID       int64     `json:"-"`
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `form:"name" json:"name" binding:"required"`
	Category string    `form:"category" json:"category" binding:"required"`
}
