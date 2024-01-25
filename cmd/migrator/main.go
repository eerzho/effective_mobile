package main

import (
	"errors"
	"flag"
	"fmt"

	"effective_mobile/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()

	m, err := migrate.New("file://"+cfg.Postgres.MigrationsPath, cfg.Postgres.Url)
	if err != nil {
		panic(err)
	}

	var down bool
	flag.BoolVar(&down, "down", false, "Set this flag to run down migrations")
	flag.Parse()

	if down {
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")

				return
			}

			panic(err)
		}

		fmt.Println("migrations reverted")
	} else {
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")

				return
			}

			panic(err)
		}

		fmt.Println("migrations applied")
	}
}
