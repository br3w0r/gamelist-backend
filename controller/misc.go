package controller

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	errType reflect.Type = errorType()
)

func ErrorSender(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

func NotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Not Found",
	})
}

func ResponseOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func ResponseInternalError(ctx *gin.Context, s string) {
	log.Printf("[ERROR] %s", s)
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal Server Error",
	})
}

func errorType() reflect.Type {
	var err error
	return reflect.ValueOf(&err).Elem().Type()
}

// Generic POST function.
// <obj> of type *T - pointer to object to be saved;
// <f> of type func(T) error - standard service save function
func GenericPost(ctx *gin.Context, obj interface{}, f interface{}) {
	// Check if f is function
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		ResponseInternalError(ctx, fmt.Sprintf("<GenericPost>: the kind of <f> is %v when must be reflect.Func", fv.Kind()))
		return
	}

	// obj's value and obj's and f's types initialization
	objVal := reflect.ValueOf(obj).Elem()
	objType := objVal.Type()
	fType := reflect.TypeOf(f)

	// Check if f's type equals standard service save function
	standardFunc := reflect.FuncOf([]reflect.Type{objType}, []reflect.Type{errType}, false)
	if fType != standardFunc {
		ResponseInternalError(ctx, fmt.Sprintf("<GenericPost>: the type of <f> is %v while it must be %v", fType, standardFunc))
		return
	}

	// Do the job
	err := ctx.ShouldBindJSON(obj)
	if err != nil {
		ErrorSender(ctx, err)
	} else {
		eVal := fv.Call([]reflect.Value{objVal})[0].Interface()
		if eVal != nil {
			ErrorSender(ctx, eVal.(error))
		} else {
			ResponseOK(ctx)
		}
	}
}
