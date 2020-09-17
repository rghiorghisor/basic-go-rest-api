package storage

import (
	"context"

	"github.com/rghiorghisor/basic-go-rest-api/model"
)

// Repository interface defining the functionality of a basic implementations.
type Repository interface {
	Create(ctx context.Context, property *model.PropertySet) error

	ReadAll(ctx context.Context) ([]*model.PropertySet, error)

	FindByID(context context.Context, id string) (*model.PropertySet, error)

	Delete(context context.Context, id string) error

	Update(ctx context.Context, property *model.PropertySet) error
}
