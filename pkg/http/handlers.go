package http

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
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

func GetFiles(r *http.Request, maxMemory int64, fieldNames ...string) (map[string]*multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, err
	}

	files := make(map[string]*multipart.FileHeader)

	for _, name := range fieldNames {
		fileList, ok := r.MultipartForm.File[name]
		if !ok || len(fileList) == 0 {
			continue
		}
		files[name] = fileList[0]
	}

	return files, nil
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
