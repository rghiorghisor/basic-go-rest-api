package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/propertyset"
	"github.com/rghiorghisor/basic-go-rest-api/server"
)

// Controller that handles the relation between the server and the service.
type Controller struct {
	service propertyset.Service
}

// PropertySetDto defines how a property set must be exposed.
type PropertySetDto struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// New retrieves a brand new contoller wrapping around the given service.
func New(service propertyset.Service) server.ControllerWrapper {
	return server.ControllerWrapper{
		Controller: &Controller{
			service: service,
		},
	}
}

type createDto struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// Create retrieves creates (if possible) a brand new property set.
func (ctrl *Controller) Create(ctx *gin.Context) {
	// Read input (must be JSON valid)
	dto := new(createDto)
	if err := ctx.BindJSON(dto); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	prop := &model.PropertySet{
		Name:   dto.Name,
		Values: dto.Values,
	}

	// Call service (business logic).
	err := ctrl.service.Create(ctx.Request.Context(), prop)

	// Respond with either error either success.
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Writer.Header().Set("Location", ctx.Request.URL.Path+"/"+prop.Name)
	ctx.Status(http.StatusCreated)
}

type readAllResponseDto struct {
	PropertySetDto []*PropertySetDto `json:"sets"`
}

// Read reads a single property set based on the provided identifier.
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
		PropertySetDto: toProperties(properties),
	})
}

type updateDto struct {
	Values []string `json:"values"`
}

// Update a single property set.
func (ctrl *Controller) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	inp := new(updateDto)
	if err := ctx.BindJSON(inp); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	prop := &model.PropertySet{
		Name:   id,
		Values: inp.Values,
	}

	if err := ctrl.service.Update(ctx.Request.Context(), prop); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, toProperty(prop))
}

// Delete a single property set, specified by means of its identifier.
func (ctrl *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := ctrl.service.Delete(ctx.Request.Context(), id); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func toProperties(bs []*model.PropertySet) []*PropertySetDto {
	out := make([]*PropertySetDto, len(bs))

	for i, b := range bs {
		out[i] = toProperty(b)
	}

	return out
}

func toProperty(b *model.PropertySet) *PropertySetDto {
	return &PropertySetDto{
		Name:   b.Name,
		Values: b.Values,
	}
}

// Register this controller to the provided group.
func (ctrl *Controller) Register(routerGroup *gin.RouterGroup) {
	api := routerGroup.Group("/set")

	api.POST("", ctrl.Create)
	api.GET("", ctrl.ReadAll)
	api.GET("/:id", ctrl.Read)
	api.PUT("/:id", ctrl.Update)
	api.DELETE("/:id", ctrl.Delete)
}
