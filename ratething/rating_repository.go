package thing

import "github.com/google/uuid"

type RatingRepository interface {
	Migrate() error
	Create(rating Rating) (*Rating, error)
	All([]Rating, error)
	GetByUUID(uuid string) (*Rating, error)
	Update(uuid uuid.UUID, updated Rating) (*Rating, error)
	Delete(uuid uuid.UUID) error
}
