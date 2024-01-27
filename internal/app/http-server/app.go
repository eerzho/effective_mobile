package http_server

import (
	"log/slog"

	"effective_mobile/internal/config"
	"effective_mobile/internal/conn"
	userH "effective_mobile/internal/http/handler/user"
	mwLogger "effective_mobile/internal/http/middleware"
	agerUserR "effective_mobile/internal/repository/agify/user"
	genderUserR "effective_mobile/internal/repository/genderize/user"
	countryerUserR "effective_mobile/internal/repository/nationalize/user"
	pUserR "effective_mobile/internal/repository/postgres/user"
	userS "effective_mobile/internal/service/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	Router *chi.Mux
	Conn   *conn.Conn
}

func New(cfg *config.Config, log *slog.Logger, connection *conn.Conn) *App {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	pUserRepo := pUserR.New(connection.Psql)
	agerUserRepo := agerUserR.New(cfg.ApiForAge)
	genderUserRepo := genderUserR.New(cfg.ApiForGen)
	countryerUserRepo := countryerUserR.New(cfg.ApiForNat)

	userService := userS.New(log, pUserRepo, agerUserRepo, genderUserRepo, countryerUserRepo)
	userHandler := userH.New(log, userService)

	router.Route("/api", func(r chi.Router) {
		r.Get("/users", userHandler.Index())
		r.Post("/users", userHandler.Store())
		r.Get("/users/{id}", userHandler.Show())
		r.Patch("/users/{id}", userHandler.Update())
		r.Delete("/users/{id}", userHandler.Delete())
	})

	router.Mount("/swagger", httpSwagger.WrapHandler)

	return &App{
		Router: router,
		Conn:   connection,
	}
}
