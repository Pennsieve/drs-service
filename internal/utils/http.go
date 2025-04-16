// internal/utils/http.go
package utils

import (
	"encoding/json"
	"net/http"

	"github.com/pennsieve/drs-service/internal/models"
)

// WriteJSON 写入JSON响应
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError 写入错误响应
func WriteError(w http.ResponseWriter, statusCode int, message string, err error) {
	// 如果有错误，将其添加到消息中
	if err != nil {
		message = message + ": " + err.Error()
	}
	
	// 根据状态码确定错误代码
	var code models.ErrorCode
	switch statusCode {
	case http.StatusBadRequest:
		code = models.ErrBadRequest
	case http.StatusUnauthorized:
		code = models.ErrUnauthorized
	case http.StatusForbidden:
		code = models.ErrForbidden
	case http.StatusNotFound:
		code = models.ErrNotFound
	case http.StatusMethodNotAllowed:
		code = models.ErrMethodNotAllowed
	case http.StatusRequestEntityTooLarge:
		code = models.ErrRequestTooLarge
	default:
		code = models.ErrInternalServer
	}
	
	errorObj := models.NewError(code, message)
	WriteJSON(w, statusCode, errorObj)
}

// WriteSpecificError 写入特定错误响应
func WriteSpecificError(w http.ResponseWriter, errorCode models.ErrorCode, message string, err error) {
	// 如果有错误，将其添加到消息中
	if err != nil {
		message = message + ": " + err.Error()
	}
	
	errorObj := models.NewError(errorCode, message)
	WriteJSON(w, errorObj.StatusCode, errorObj)
}
