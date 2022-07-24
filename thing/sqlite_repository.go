package thing

import "github.com/google/uuid"

type Repository interface {
	Migrate() error
	Create(thing Thing) (*Thing, error)
	All([]Thing, error)
	GetByName(name string) (*Thing, error)
	Update(uuid uuid.UUID, updated Thing) (*Thing, error)
	Delete(uuid uuid.UUID) error
}
