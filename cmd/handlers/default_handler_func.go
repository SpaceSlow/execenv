package handlers

import "net/http"

func DefaultHandlerFunc(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}
