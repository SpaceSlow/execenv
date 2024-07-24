package handlers

import "net/http"

// BadRequestHandlerFunc вовзращает 400 Bad Request.
func BadRequestHandlerFunc(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}
