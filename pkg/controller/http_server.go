package controller

import (
	"context"
	"net/http"
)

func listenAndServeHttp(ctx context.Context, handler http.Handler) error {
	return http.ListenAndServe(":80", handler)
}
