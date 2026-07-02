package httputil

import (
	"FeatureFlags/internal/apperror"
	"FeatureFlags/internal/validation"
	"encoding/json"
	"errors"
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

// ==========================================
// ЗАПИСЬ ОТВЕТОВ (Response)
// ==========================================

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(status)

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func HandleError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError

	if errors.As(err, &appErr) {
		WriteJSON(w, appErr.Code, map[string]string{"error": appErr.Message})
		return
	}

	WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "внутренняя ошибка сервера"})
}
