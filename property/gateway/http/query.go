package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rghiorghisor/basic-go-rest-api/property"
)

func parse(ctx *gin.Context) property.Query {
	q := property.Query{}

	q.ID = ctx.Param("id")
	q.Set = ctx.Query("set")
	q.Fields = property.NewFields(ctx.QueryArray("fields"))

	return q
}
