package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"chat_app/internal/config"
	"chat_app/internal/migrations"
)

func main() {
	var (
		action  = flag.String("action", "up", "Migration action: up, down, status")
		steps   = flag.Int("steps", 0, "Number of steps to migrate (0 = all)")
		envFile = flag.String("env", "", "Environment file to load (e.g., .env, env.dev)")
	)
	flag.Parse()

	var cfg *config.Config
	if *envFile != "" {
		cfg = config.LoadFromFile(*envFile)
	} else {
		cfg = config.Load()
	}

	db, err := config.NewDatabaseConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.CloseDatabase(db)

	switch *action {
	case "up":
		fmt.Println("Running migrations...")
		if err := migrations.RunMigrations(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully!")

	case "status":
		version, err := migrations.GetCurrentVersion(db)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
		fmt.Printf("Current migration version: %d\n", version)

	case "down":
		if *steps <= 0 {
			log.Fatal("Steps must be specified for down migration (e.g., -steps=1)")
		}
		fmt.Printf("Rolling back %d migration(s)...\n", *steps)
		if err := migrations.RollbackMigrations(db, *steps); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("Rollback completed successfully!")

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, down, status")
		os.Exit(1)
	}
}
