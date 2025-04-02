package ddd

import "context"

type Void struct{}

// Repository defines methods for a repository that stores entities
type Repository[T any] interface {
	GetNextId(ctx context.Context) (uint, error)
	Store(ctx context.Context, entity *T) error  // Stores the entity in the repository
}

// ApplicationService represents a service that processes a payload and returns a result
type ApplicationService[P, R any] func(ctx context.Context, payload P) (R, error)

// DataMapper defines methods to map data to and from an entity
type DataMapper[ENTITY any, DATA any] interface {
	ToEntity(data DATA) ENTITY // Converts data to an entity
	ToData(entity ENTITY) DATA // Converts an entity to data
}
