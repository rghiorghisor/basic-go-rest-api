package storage

import (
	"context"

	"github.com/rghiorghisor/basic-go-rest-api/model"
)

// Repository interface defining the functionality of a basic implementations.
type Repository interface {
	Create(ctx context.Context, property *model.Property) error

	ReadAll(ctx context.Context) ([]*model.Property, error)

	FindByID(context context.Context, id string) (*model.Property, error)

	FindByName(context context.Context, name string) (*model.Property, error)

	Delete(context context.Context, id string) error

	Update(ctx context.Context, property *model.Property) error
}
