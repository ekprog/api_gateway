package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GeneralMW(ctx *gin.Context) {
	ctx.Header("content-type", "application/json")
}

func ErrorMW(ctx *gin.Context) {
	ctx.Next()
	if len(ctx.Errors) > 0 {
		err := ctx.Errors[0].Err

		switch err.(type) {
		case *json.UnmarshalTypeError:
			e := err.(*json.UnmarshalTypeError)
			ctx.JSON(500, ValidationError(fmt.Sprintf("%s (type) = %s", e.Field, e.Type)))
		case validator.ValidationErrors:
			errs := err.(validator.ValidationErrors)
			if len(errs) > 0 {
				ctx.JSON(500, ValidationError(fmt.Sprintf("%s -> %s", errs[0].Field(), errs[0].ActualTag())))
			} else {
				ctx.JSON(500, ValidationError())
			}
		default:
			ctx.JSON(500, ServerError())
		}
	}
}
