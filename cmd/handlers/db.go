package handlers

import (
	"context"
	"github.com/jackc/pgx/v5"
	"net/http"
)

type DBHandler struct {
	Ctx  context.Context
	Conn *pgx.Conn
}

func (h DBHandler) Ping(res http.ResponseWriter, _ *http.Request) {
	if h.Conn != nil && h.Conn.Ping(h.Ctx) == nil {
		res.WriteHeader(http.StatusOK)
		return
	}
	res.WriteHeader(http.StatusInternalServerError)
}
