package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/property"
)

// Controller that handles the relation between the server and the service.
type Controller struct {
	service property.Service
}

// PropertyDto defines how a property must be exposed.
type PropertyDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// NewController retrieves a brand new contoller wrapping around the given service.
func NewController(service property.Service) *Controller {
	return &Controller{
		service: service,
	}
}

type createDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// Create retrieves creates (if possible) a brand new property.
func (ctrl *Controller) Create(ctx *gin.Context) {
	// Read input (must be JSON valid)
	dto := new(createDto)
	if err := ctx.BindJSON(dto); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	prop := &model.Property{
		Name:        dto.Name,
		Description: dto.Description,
		Value:       dto.Value,
	}

	// Call service (business logic).
	err := ctrl.service.Create(ctx.Request.Context(), prop)

	// Respond with either error either success.
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Writer.Header().Set("Location", ctx.Request.URL.Path+"/"+prop.ID)
	ctx.Status(http.StatusCreated)
}

type readAllResponseDto struct {
	PropertyDto []*PropertyDto `json:"properties"`
}

func (ctrl *Controller) Read(ctx *gin.Context) {
	id := ctx.Param("id")

	foundProp, err := ctrl.service.FindByID(ctx.Request.Context(), id)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, toProperty(foundProp))
}

// ReadAll retrieves a list of all available properties.
func (ctrl *Controller) ReadAll(ctx *gin.Context) {
	properties, err := ctrl.service.ReadAll(ctx.Request.Context())

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, &readAllResponseDto{
		PropertyDto: toProperties(properties),
	})
}

type updateDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// Update a single property.
func (ctrl *Controller) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	inp := new(updateDto)
	if err := ctx.BindJSON(inp); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	prop := &model.Property{
		ID:          id,
		Name:        inp.Name,
		Description: inp.Description,
		Value:       inp.Value,
	}

	err := ctrl.service.Update(ctx.Request.Context(), prop)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, toProperty(prop))
}

// Delete a single property, specified by means of its identifier.
func (ctrl *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := ctrl.service.Delete(ctx.Request.Context(), id)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func toProperties(bs []*model.Property) []*PropertyDto {
	out := make([]*PropertyDto, len(bs))

	for i, b := range bs {
		out[i] = toProperty(b)
	}

	return out
}

func toProperty(b *model.Property) *PropertyDto {
	return &PropertyDto{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		Value:       b.Value,
	}
}
