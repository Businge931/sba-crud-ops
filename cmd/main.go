package main

import (
	"github.com/Businge931/api-gateway/internal/app/handlers"
	"github.com/Businge931/api-gateway/internal/db"
	"github.com/Businge931/api-gateway/internal/env"
	"github.com/Businge931/api-gateway/internal/repository/postgres"
	"github.com/Businge931/api-gateway/internal/service"

	log "github.com/sirupsen/logrus"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://admin:adminpassword@localhost:5432/api_gateway?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_OPEN_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	log.Printf("Using database connection string: %s", cfg.db.addr)

	pool, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer pool.Close()
	log.Print("database connection pool established")

	// Initialize repository
	oddsRepo := postgres.NewOddsRepository(pool)

	// Initialize service
	oddsService := service.NewOddsService(oddsRepo)

	// Initialize handlers
	oddsHandler := handlers.NewOddsHandler(oddsService)

	app := &application{
		config:      cfg,
		oddsHandler: oddsHandler,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
