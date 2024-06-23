package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	db "go-app/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var fs embed.FS

func main() {
	// Env
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	listenAddr := os.Getenv("LISTEN_ADDR")
	isProd := os.Getenv("PROD") == "1"

	dbPath := os.Getenv("DB_PATH")

	// Logging
	loggerLevel := slog.LevelDebug
	if isProd {
		loggerLevel = slog.LevelInfo
	}
	logger := httplog.NewLogger("biogutterne-fastapi", httplog.Options{
		LogLevel:         loggerLevel,
		Concise:          true,
		RequestHeaders:   false,
		MessageFieldName: "message",
		TimeFieldFormat:  time.RFC1123,
		Tags:             map[string]string{},
		QuietDownRoutes: []string{
			"/",
			"/ping",
			"/favicon.ico",
		},
		QuietDownPeriod: 20 * time.Second,
	})

	logger.Info(fmt.Sprintf("Starting application with LogLevel=%s IsProd=%v", loggerLevel, isProd))

	// Database
	migrateDB := func() {
		sqliteDB, err := sql.Open("sqlite3", fmt.Sprintf("file:%v", dbPath))
		if err != nil {
			panic(err)
		}
		defer sqliteDB.Close()
		migrations, err := iofs.New(fs, "migrations")
		if err != nil {
			panic(err)
		}
		logger.Info("Database OK, applying migrations")
		driver, err := sqlite3.WithInstance(sqliteDB, &sqlite3.Config{})
		if err != nil {
			panic(err)
		}

		m, err := migrate.NewWithInstance(
			"iofs", migrations,
			"sqlite3", driver)
		if err != nil {
			panic(err)
		}

		err = m.Up()
		if err != nil {
			logger.Info("No migrations", "result", err)
		} else {
			logger.Info("Migrations OK")
		}
	}
	migrateDB()

	// Re-open DB for queries
	sqliteDB, err := sql.Open("sqlite3", fmt.Sprintf("file:%v", dbPath))
	if err != nil {
		panic(err)
	}
	defer sqliteDB.Close()
	queries := db.New(sqliteDB)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.GetHead)
	r.Use(func(next http.Handler) http.Handler {
		// Logging middleware
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now().UTC()
			defer (func() {
				st := http.StatusText(ww.Status())
				var level slog.Level
				if ww.Status() <= 0 || ww.Status() >= 500 {
					level = slog.LevelWarn
				} else {
					level = slog.LevelDebug
				}
				logger.Log(r.Context(), level, fmt.Sprintf("%v %v %v - %v %v", r.Method, ww.Status(), st, time.Since(t1), r.URL.Path))
			})()

			next.ServeHTTP(ww, r)
		})
	})
	r.Use(loggerMiddleware(logger.Logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Hello World!\n")
	})

	r.Get("/cities", func(w http.ResponseWriter, r *http.Request) {
		log := useLogger(r)

		log.Debug("Fetching cities")
		ret, err := queries.GetAllCities(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		render.JSON(w, r, ret)
	})

	// Run
	logger.Info(fmt.Sprintf("Listening on %v", listenAddr))
	err = http.ListenAndServe(listenAddr, r)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen: %v", err))
	}
}

// -- Middleware --
// ----------------

// --- Middleware helpers
type hKey int

const (
	hLoggerKey hKey = iota
)

func useLogger(r *http.Request) *slog.Logger {
	logger, ok := (r.Context().Value(hLoggerKey)).(*slog.Logger)
	if !ok || logger == nil {
		panic("request logger not defined")
	}
	return logger
}

// --- HTTP Middleware
func loggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), hLoggerKey, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
