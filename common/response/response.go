package response

import (
	"context"
	"net/http"
)

type JsonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OkHandler(ctx context.Context, v any) any {
	return JsonResponse{Code: http.StatusOK, Message: "ok", Data: v}
}
