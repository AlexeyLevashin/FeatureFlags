package apperror

import (
	"encoding/json"
	"errors"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	var appErr *AppError
	w.Header().Set("Content-Type", "application/json")
	if errors.As(err, &appErr) {
		w.WriteHeader(appErr.Code)
		json.NewEncoder(w).Encode(map[string]string{"error": appErr.Message})
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": "внутренняя ошибка сервера"})
}
