package handlers

import (
	"FeatureFlags/internal/apperror"
	"FeatureFlags/internal/validation"
	"encoding/json"
	"net/http"
)

func ReadAndValidate[T any](r *http.Request) (T, error) {
	var payload T

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload, apperror.BadRequest("Неверный формат JSON")
	}

	if err := validation.Validate.Struct(payload); err != nil {
		return payload, apperror.BadRequest("Некорректные данные")
	}

	return payload, nil
}
