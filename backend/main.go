package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mosanda-vacation-attendace/handlers"
	"mosanda-vacation-attendace/router"
	"mosanda-vacation-attendace/service"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // source driver
	_ "github.com/lib/pq"                                      // Blank import loads the driver
)

func runMigrations(parConnStr string) error {
	m, pvError := migrate.New(
		"file://database/migrations", // path to your SQL files folder
		parConnStr,
	)
	if pvError != nil {
		return fmt.Errorf("MIG0001 %v", pvError)
	}

	// Run migrations
	if pvError := m.Up(); pvError != nil && pvError != migrate.ErrNoChange {
		return fmt.Errorf("MIG0002 %v", pvError)
	}

	log.Println("Migrations applied successfully!")
	return nil
}

func main() {
	connStr := "postgres://postgres:admin@localhost:5432/mosanda?sslmode=disable"

	// 1. Run migrations
	if err := runMigrations(connStr); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 2. Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// 3. Verify database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	// 4. Seed initial users
	if err := seedUsers(db); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	// 5. Start application
	appService := service.NewService(nil, db)
	startExpiredSessionCleanup(appService)

	appHandlers := handlers.NewHandlers(appService)
	appRouter := router.NewRouter(appHandlers)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", appRouter); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func startExpiredSessionCleanup(parService *service.Service) {
	go func() {
		pvTicker := time.NewTicker(15 * time.Minute)
		defer pvTicker.Stop()

		for {
			pvContext, pvCancel := context.WithTimeout(context.Background(), 10*time.Second)
			pvDeletedCount, pvError := parService.DeleteExpiredSessions(pvContext)
			pvCancel()

			if pvError != nil {
				log.Printf("Expired session cleanup failed: %v", pvError)
			} else if pvDeletedCount > 0 {
				log.Printf("Deleted %d expired sessions", pvDeletedCount)
			}

			<-pvTicker.C
		}
	}()
}

func seedUsers(db *sql.DB) error {
	users := map[string]string{
		"ali@ali.com":           "ali",
		"hisham@hisham.com":     "hisham",
		"abdo@abdo.com":         "abdo",
		"abdallah@abdallah.com": "abdallah",
		"amir@amir.com":         "amir",
	}

	ctx := context.Background()

	for email, password := range users {
		var exists bool

		err := db.QueryRowContext(
			ctx,
			"SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)",
			email,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("check user %s: %w", email, err)
		}

		if exists {
			continue
		}

		passwordHash, err := argon2id.CreateHash(
			password,
			argon2id.DefaultParams,
		)
		if err != nil {
			return fmt.Errorf("hash password for %s: %w", email, err)
		}

		_, err = db.ExecContext(
			ctx,
			`INSERT INTO users (email, hash_password, created_at)
			 VALUES ($1, $2, $3)
			`,
			email,
			passwordHash,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("seed user %s: %w", email, err)
		}
	}

	log.Println("User seeding completed")
	return nil
}
