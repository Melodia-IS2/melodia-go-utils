package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"

	pkgErrors "github.com/Melodia-IS2/melodia-go-utils/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func JSON(w http.ResponseWriter, status int, body any) {
	if body != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	} else {
		w.WriteHeader(status)
	}
}

func Ok(w http.ResponseWriter, body any) {
	JSON(w, http.StatusOK, body)
}

func Created(w http.ResponseWriter, body any) {
	JSON(w, http.StatusCreated, body)
}

func NoContent(w http.ResponseWriter) {
	JSON(w, http.StatusNoContent, nil)
}

func GetUrlParam[T any](r *http.Request, param string) (T, error) {
	var value T
	paramValue := chi.URLParam(r, param)
	if paramValue == "" {
		return value, errors.New("url param not found")
	}
	result, err := parseParam[T](paramValue)
	if err != nil {
		return result, pkgErrors.NewBadRequestError(err.Error())
	}
	return result, nil
}

func GetQueryParam[T any](r *http.Request, param string) (*T, error) {
	paramValue := r.URL.Query().Get(param)
	if paramValue == "" {
		return nil, nil
	}
	result, err := parseParam[T](paramValue)
	if err != nil {
		return &result, pkgErrors.NewBadRequestError(err.Error())
	}
	return &result, nil
}

func parseParam[T any](paramValue string) (T, error) {
	var value T
	var anyValue any

	switch any(value).(type) {
	case string:
		anyValue = paramValue
	case int:
		v, err := strconv.Atoi(paramValue)
		if err != nil {
			return value, fmt.Errorf("invalid value: %s. Expected int", paramValue)
		}
		anyValue = v
	case uuid.UUID:
		v, err := uuid.Parse(paramValue)
		if err != nil {
			return value, fmt.Errorf("invalid value: %s. Expected uuid", paramValue)
		}
		anyValue = v
	case bool:
		v, err := strconv.ParseBool(paramValue)
		if err != nil {
			return value, fmt.Errorf("invalid value: %s. Expected bool", paramValue)
		}
		anyValue = v
	default:
		return value, fmt.Errorf("unsupported type: %T", value)
	}

	return anyValue.(T), nil
}

func MapRequest[T any](r *http.Request) (T, error) {
	var request T
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}
	return request, nil
}

func MapFormRequest[T any](r *http.Request) (T, error) {
	var request T
	v := reflect.ValueOf(&request).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		formTag := field.Tag.Get("form")

		if formTag == "" {
			continue
		}

		if field.Type == reflect.TypeOf(&multipart.FileHeader{}) {
			if fileHeaders, ok := r.MultipartForm.File[formTag]; ok && len(fileHeaders) > 0 {
				v.Field(i).Set(reflect.ValueOf(fileHeaders[0]))
			}
			continue
		}

		formValue := r.FormValue(formTag)
		if formValue == "" {
			continue
		}

		fieldVal := v.Field(i)

		switch field.Type.Kind() {
		case reflect.String:
			fieldVal.SetString(formValue)

		case reflect.Slice, reflect.Struct:
			newVal := reflect.New(field.Type)
			if err := json.Unmarshal([]byte(formValue), newVal.Interface()); err != nil {
				return request, fmt.Errorf("failed to unmarshal JSON for field %s: %w", field.Name, err)
			}
			fieldVal.Set(newVal.Elem())

		default:
		}
	}

	return request, nil
}
