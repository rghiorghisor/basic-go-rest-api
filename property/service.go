package property

import (
	"context"

	"github.com/rghiorghisor/basic-go-rest-api/model"
)

// Service defines the use case available for properties
type Service interface {
	Create(ctx context.Context, property *model.Property) error

	ReadAll(ctx context.Context, query Query) ([]*model.Property, error)

	FindByID(ctx context.Context, id string) (*model.Property, error)

	Delete(ctx context.Context, id string) error

	Update(ctx context.Context, property *model.Property) error
}
