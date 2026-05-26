package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // registers the postgres driver with database/sql; blank import is intentional
	"github.com/pressly/goose/v3"
	"github.com/FatihTirek/football-league/migrations"

	"github.com/FatihTirek/football-league/internal/engine"
	"github.com/FatihTirek/football-league/internal/handler"
	"github.com/FatihTirek/football-league/internal/repository"
	"github.com/FatihTirek/football-league/internal/service"
)

func main() {
	// ── 1. DATABASE ───────────────────────────────────────────────────────────
	dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSLMODE")

	if dbUser == "" || dbPass == "" || dbName == "" || dbHost == "" || dbPort == "" || sslMode == "" {
        log.Fatalf("one or more required database environment variables are missing.")
    }

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost,
		dbPort,
		dbUser,
		dbPass,
		dbName,
		sslMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// 2. Tell Goose to read your 001_initial_schema.sql
    goose.SetBaseFS(migrations.FS)

	// 3. Run it!
    if err := goose.Up(db, "."); err != nil {
        log.Fatal("Failed to run migrations: ", err)
    }

    log.Println("Database migrated and ready!")

	// Tune the connection pool. For this 4-player simulation scale these
	// defaults are generous, but they are good habits to set explicitly.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the database is actually reachable before accepting HTTP traffic.
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer pingCancel()
	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("Failed to reach database: %v", err)
	}
	log.Println("database connection established")

	// ── 2. DEPENDENCY WIRING (bottom-up) ─────────────────────────────────────
	// Repositories depend only on *sql.DB.
	teamRepository     := repository.NewTeamRepository(db)
	matchRepository    := repository.NewMatchRepository(db)
	standingRepository := repository.NewStandingRepository(db)

	// Engines are pure math — no DB, no HTTP.
	matchEngine      := engine.NewMatchEngine(time.Now().UnixNano())
	predictionEngine := engine.NewPredictionEngine()

	// Services depend on repository interfaces and engines.
	leagueService     := service.NewLeagueService(teamRepository, matchRepository, standingRepository, matchEngine)
	predictionService := service.NewPredictionService(teamRepository, matchRepository, standingRepository, predictionEngine)

	// The router is the topmost layer — it depends only on service interfaces.
	router := handler.NewRouter(leagueService, predictionService)

	// ── 3. HTTP SERVER ────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      router,
		ReadTimeout:  15 * time.Second, // max time to read request headers + body
		WriteTimeout: 30 * time.Second, // max time to write the full response (generous for Monte Carlo)
		IdleTimeout:  60 * time.Second, // max time a keep-alive connection stays open
	}

	// Start the server in a goroutine so it doesn't block the shutdown logic below.
	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// ── 4. GRACEFUL SHUTDOWN ──────────────────────────────────────────────────
	// Block here until the OS sends SIGINT (Ctrl+C) or SIGTERM (Docker stop).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown signal received, draining active requests...")

	// Give in-flight requests up to 10 seconds to finish before force-closing.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced shutdown: %v", err)
	}
	log.Println("server stopped cleanly")
}