package thing

import (
	"github.com/google/uuid"
)

type Rating struct {
	ID         int64     `json:"-"`
	UUID       uuid.UUID `json:"uuid"`
	Candidates string    `json:"-"`
	Winner     uuid.UUID `json:"winner"`
	Loser      uuid.UUID `json:"loser"`
}
