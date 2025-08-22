package basehandler

import "github.com/gin-gonic/gin"

type TScope[T any] func(*T) error
type TWithIDScope[T any] func(*T, uint64) error
type TAfterValidate[T any] func(ctx *gin.Context, modelDTO *T, fields *[]string) error
type TResponseHook[T any] func(ctx *gin.Context, model *T) error
type TResponseListHook[T any] func(ctx *gin.Context, rows *[]T) error
