package http

import (
	"bytes"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
	"github.com/rghiorghisor/basic-go-rest-api/model"
)

type formatters struct {
	values []formatter
}

func newFormatters() formatters {
	return formatters{
		values: []formatter{
			&jsonFormatter{},
			&javaPropertiesFormatter{},
		},
	}
}

func (f formatters) process(ctx *gin.Context, code int, bs []*model.Property) {
	acceptHeader := ctx.Request.Header.Get("Accept")
	found := false
	for _, f := range f.values {
		if !f.supports(acceptHeader) {
			continue
		}

		f.process(ctx, code, bs)
		found = true
		break
	}

	if !found {
		f.values[0].process(ctx, code, bs)
	}
}

type formatter interface {
	supports(acceptHeader string) bool
	process(ctx *gin.Context, code int, bs []*model.Property)
}

type jsonFormatter struct {
}

func (f jsonFormatter) supports(acceptHeader string) bool {
	return acceptHeader == "" ||
		strings.EqualFold(acceptHeader, "*/*") ||
		strings.EqualFold(acceptHeader, "application/json")
}

func (f jsonFormatter) process(ctx *gin.Context, code int, bs []*model.Property) {
	ctx.JSON(code, &readAllResponseDto{
		PropertyDto: toProperties(bs),
	})
}

type javaPropertiesFormatter struct {
}

func (f javaPropertiesFormatter) supports(acceptHeader string) bool {
	return strings.EqualFold(acceptHeader, "application/java.properties")
}

func (f javaPropertiesFormatter) process(ctx *gin.Context, code int, bs []*model.Property) {

	buf := new(bytes.Buffer)
	props := properties.NewProperties()
	for _, prop := range bs {
		props.Set(prop.Name, prop.Value)
		if prop.Description != "" {
			props.SetComment(prop.Name, prop.Description)
		}
	}

	props.WriteComment(buf, "# ", properties.UTF8)

	ctx.Header("Content-Disposition", `attachment; filename="java.properties"`)
	ctx.Data(code, "application/octet-stream", buf.Bytes())
}
