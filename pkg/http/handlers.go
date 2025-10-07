package http

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	pkgErrors "github.com/Melodia-IS2/melodia-go-utils/pkg/errors"
	"github.com/Melodia-IS2/melodia-go-utils/pkg/router"

	"github.com/go-playground/validator"
)

var exposeErrorDetail = false

func SetExposeErrorDetail(value bool) {
	exposeErrorDetail = value
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var appError *pkgErrors.AppError
	if errors.As(err, &appError) {
		router.JSON(w, appError.HTTPCode, pkgErrors.ErrorResponse{
			Type:     "about:blank",
			Title:    appError.Title,
			Status:   appError.HTTPCode,
			Detail:   appError.Message,
			Instance: r.URL.Path,
		})
	} else {
		detail := "An unexpected error occurred"
		if exposeErrorDetail {
			detail = err.Error()
		}
		router.JSON(w, http.StatusInternalServerError, pkgErrors.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   detail,
			Instance: r.URL.Path,
		})
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	router.JSON(w, http.StatusNotFound, pkgErrors.ErrorResponse{
		Type:     "about:blank",
		Title:    "Not Found",
		Status:   http.StatusNotFound,
		Detail:   "The requested resource was not found",
		Instance: r.URL.Path,
	})
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	router.JSON(w, http.StatusMethodNotAllowed, pkgErrors.ErrorResponse{
		Type:     "about:blank",
		Title:    "Method Not Allowed",
		Status:   http.StatusMethodNotAllowed,
		Detail:   "The requested method is not allowed for this resource",
		Instance: r.URL.Path,
	})
}

func ParseBody[T any](r *http.Request) (T, error) {
	request, err := router.MapRequest[T](r)
	if err != nil {
		return request, pkgErrors.NewBadRequestError("Invalid request body")
	}

	if err := validator.New().Struct(request); err != nil {
		return request, pkgErrors.NewValidationError(formatValidationError(err))
	}
	return request, nil
}

func MapRequestForm[T any](r *http.Request, maxMemoryMB int64) (*T, error) {
	if err := r.ParseMultipartForm(maxMemoryMB << 20); err != nil {
		return nil, err
	}

	result := new(T)
	val := reflect.ValueOf(result).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("form")
		if tag == "" {
			continue
		}

		fv := val.Field(i)

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(r.FormValue(tag))
		case reflect.Bool:
			v := r.FormValue(tag)
			fv.SetBool(v == "true" || v == "1")
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v := r.FormValue(tag)
			if v != "" {
				num, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return nil, err
				}
				fv.SetInt(num)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v := r.FormValue(tag)
			if v != "" {
				num, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return nil, err
				}
				fv.SetUint(num)
			}
		case reflect.Ptr:
			if fv.Type() == reflect.TypeOf(&multipart.FileHeader{}) {
				file, header, err := r.FormFile(tag)
				if err == nil {
					fv.Set(reflect.ValueOf(header))
					file.Close()
				}
			} else {
				return nil, errors.New("unsupported pointer type: " + fv.Type().String())
			}
		default:
			return nil, errors.New("unsupported field type: " + fv.Type().String())
		}
	}

	return result, nil
}

func formatValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		messages := make([]string, 0, len(errs))
		for _, e := range errs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("%s is required", field))
			case "min":
				messages = append(messages, fmt.Sprintf("%s must be at least %s characters", field, e.Param()))
			case "max":
				messages = append(messages, fmt.Sprintf("%s must be less than %s characters", field, e.Param()))
			default:
				messages = append(messages, fmt.Sprintf("%s is invalid", field))
			}
		}
		return strings.Join(messages, ", ")
	}
	return err.Error()
}
