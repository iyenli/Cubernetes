package httpserver

import (
	"Cubernetes/pkg/gateway/httpserver/handlers"
	"net/http"
)

var serverlessList = []Handler{
	{http.MethodGet, "/health", handlers.Hello},
}
