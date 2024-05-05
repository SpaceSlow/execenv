package handlers

import "net/http"

func BadRequestHandlerFunc(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}
