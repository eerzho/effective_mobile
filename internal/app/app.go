package app

import (
	"log/slog"

	"effective_mobile/internal/app/http-server"
	"effective_mobile/internal/conn"
)

type App struct {
	HttpApp *http_server.App
}

func New(log *slog.Logger, connection *conn.Conn) *App {
	return &App{HttpApp: http_server.New(log, connection)}
}
