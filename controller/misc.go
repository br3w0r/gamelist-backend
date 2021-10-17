package controller

import (
	"fmt"
	"net/http"
	"reflect"

	utilErrs "github.com/br3w0r/gamelist-backend/util/errors"
	"github.com/gin-gonic/gin"
)

var (
	errType reflect.Type = errorType()
)

const (
	errPost = "failed to process post request"
)

func ErrorSender(ctx *gin.Context, err error) {
	var utilErr *utilErrs.Error
	if castErr, ok := err.(*utilErrs.Error); ok {
		utilErr = castErr
	} else {
		utilErr = utilErrs.New(utilErrs.Internal, err, "unknown error")
	}

	ctx.AbortWithStatusJSON(utilErr.Code().ToHTTP(), utilErr)
}

func ResponseOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
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
		err := fmt.Errorf("the kind of <f> is %v when must be reflect.Func", fv.Kind())
		ErrorSender(ctx, utilErrs.New(utilErrs.Internal, err, errPost))
		return
	}

	// obj's value and obj's and f's types initialization
	objVal := reflect.ValueOf(obj).Elem()
	objType := objVal.Type()
	fType := reflect.TypeOf(f)

	// Check if f's type equals standard service save function
	standardFunc := reflect.FuncOf([]reflect.Type{objType}, []reflect.Type{errType}, false)
	if fType != standardFunc {
		err := fmt.Errorf("<GenericPost>: the type of <f> is %v while it must be %v", fType, standardFunc)
		ErrorSender(ctx, utilErrs.New(utilErrs.Internal, err, errPost))
		return
	}

	// Do the job
	err := ctx.ShouldBindJSON(obj)
	if err != nil {
		ErrorSender(ctx, utilErrs.JSONParseErr(err))
		return
	}

	eVal := fv.Call([]reflect.Value{objVal})[0].Interface()
	if eVal != nil {
		ErrorSender(ctx, eVal.(error))
		return
	}

	ResponseOK(ctx)
}
