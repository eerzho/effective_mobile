package conn

import (
	"effective_mobile/internal/config"
	"effective_mobile/internal/conn/postgres"
)

type Conn struct {
	Psql *postgres.Conn
}

func MustNew(cfg *config.Config) *Conn {
	return &Conn{Psql: postgres.MustNew(cfg.Postgres.Url)}
}
