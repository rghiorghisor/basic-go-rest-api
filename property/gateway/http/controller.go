package http

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/property"
	"github.com/rghiorghisor/basic-go-rest-api/server"
)

// Controller that handles the relation between the server and the service.
type Controller struct {
	formatters formatters
	service    property.Service
}

// PropertyDto defines how a property must be exposed.
type PropertyDto struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value"`
}

// New retrieves a brand new contoller wrapping around the given service.
func New(service property.Service) server.ControllerWrapper {
	return server.ControllerWrapper{
		Controller: &Controller{
			formatters: newFormatters(),
			service:    service,
		},
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
	PropertyDto []PropertyDto `json:"properties"`
}

func (ctrl *Controller) Read(ctx *gin.Context) {
	ctrl.readOne(ctx, toPropertyFiltered)
}

// ReadBasic retrieves a basic form of a property, i.e only name and value.
func (ctrl *Controller) ReadBasic(ctx *gin.Context) {
	ctrl.readOne(ctx, toBasicProperty)
}

func (ctrl *Controller) readOne(ctx *gin.Context, f func(*model.Property, property.Query) interface{}) {
	query := parse(ctx)

	foundProp, err := ctrl.service.FindByID(ctx.Request.Context(), query.ID)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, f(foundProp, query))
}

// ReadAll retrieves a list of all available properties.
func (ctrl *Controller) ReadAll(ctx *gin.Context) {
	query := parse(ctx)
	properties, err := ctrl.service.ReadAll(ctx.Request.Context(), query)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctrl.formatters.process(ctx, http.StatusOK, properties)
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

	ctx.JSON(http.StatusOK, toPropertyDTO(prop))
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

func toProperties(bs []*model.Property) []PropertyDto {
	out := make([]PropertyDto, len(bs))

	for i, b := range bs {
		out[i] = toPropertyDTO(b)
	}

	return out
}

func toBasicProperty(b *model.Property, query property.Query) interface{} {
	return PropertyDto{
		Name:  b.Name,
		Value: b.Value,
	}
}

func toPropertyDTO(b *model.Property) PropertyDto {
	return PropertyDto{
		ID:          b.ID,
		Name:        b.Name,
		Value:       b.Value,
		Description: b.Description,
	}
}

func toPropertyFiltered(b *model.Property, query property.Query) interface{} {
	dto := toPropertyDTO(b)

	if !query.Fields.IsEnabled() {
		return dto
	}

	rt := reflect.TypeOf(dto)
	rv := reflect.ValueOf(dto)

	out := make(map[string]interface{}, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		jsonKey := field.Tag.Get("json")

		if query.Fields.Contains(jsonKey) {
			out[jsonKey] = rv.Field(i).Interface()
		}
	}

	return out
}

// Register this controller to the provided group.
func (ctrl *Controller) Register(routerGroup *gin.RouterGroup) {
	api := routerGroup.Group("/property")

	api.POST("", ctrl.Create)
	api.GET("", ctrl.ReadAll)
	api.GET("/:id", ctrl.Read)
	api.GET("/:id/basic", ctrl.ReadBasic)
	api.PUT("/:id", ctrl.Update)
	api.DELETE("/:id", ctrl.Delete)
}
