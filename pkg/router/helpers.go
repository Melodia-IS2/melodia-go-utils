package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

func wrapData(data any) any {
	return map[string]any{
		"data": data,
	}
}

func Ok(w http.ResponseWriter, body any) {
	JSON(w, http.StatusOK, wrapData(body))
}

func Created(w http.ResponseWriter, body any) {
	JSON(w, http.StatusCreated, wrapData(body))
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
