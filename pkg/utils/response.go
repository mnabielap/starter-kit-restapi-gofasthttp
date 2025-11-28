package utils

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// Response standardizes the JSON response format
type Response struct {
	Status  string      `json:"status"` // "success" or "error"
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// WriteJSON sends a JSON response to the fasthttp context
func WriteJSON(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	
	if err := json.NewEncoder(ctx).Encode(data); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

// WriteError sends a standard error response
func WriteError(ctx *fasthttp.RequestCtx, statusCode int, message string) {
	WriteJSON(ctx, statusCode, Response{
		Status:  "error",
		Message: message,
	})
}

// WriteSuccess sends a standard success response
func WriteSuccess(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) {
	WriteJSON(ctx, statusCode, data)
}