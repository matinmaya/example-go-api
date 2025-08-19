package basehandler

import "github.com/gin-gonic/gin"

type TScope[T any] func(*T) error
type TScopeWithID[T any] func(*T, uint64) error
type TAfterValidate[T any] func(ctx *gin.Context, modelDTO *T, fields *[]string) error
type TBeforeResponse[T any] func(ctx *gin.Context, model *T) error
