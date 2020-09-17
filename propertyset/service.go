package propertyset

import (
	"context"

	"github.com/rghiorghisor/basic-go-rest-api/model"
)

// Service defines the use case available for properties
type Service interface {
	Create(ctx context.Context, property *model.PropertySet) error

	ReadAll(ctx context.Context) ([]*model.PropertySet, error)

	FindByID(ctx context.Context, id string) (*model.PropertySet, error)

	Delete(ctx context.Context, id string) error

	Update(ctx context.Context, property *model.PropertySet) error
}
