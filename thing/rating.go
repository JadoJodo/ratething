package thing

import "github.com/google/uuid"

type Rating struct {
	id         int64
	candidates string
	winner     uuid.UUID
	loser      uuid.UUID
}
