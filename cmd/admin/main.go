// This program performs administrative tasks for the garage sale service.
package main

import (
	"context"
	"github.com/deezone/HydroBytes-BaseStation/internal/account"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/auth"
	"time"

	// Core Packages
	"fmt"
	"log" // https://golang.org/pkg/log/
	"os"

	// Third-party packages
	"github.com/pkg/errors"

	// Internal applcation packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/conf"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/database"
	"github.com/deezone/HydroBytes-BaseStation/internal/schema"
)

// Main entry point for program.
func main() {

	// Only call Exit in main() to allow all defers to complete before shutdown in the case of an error
	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

// Main application logic.
func run() error {

	// =========================================================================
	// Configuration

	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "STATIONS", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("STATIONS", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "error: parsing config")
	}

	// =========================================================================
	// Database configuration

	// Initialize dependencies.
	dbConfig := database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	}

	// =========================================================================
	// Supported admin commands

	var err error
	switch cfg.Args.Num(0) {
	case "migrate":
		err = migrate(dbConfig)
	case "seed":
		err = seed(dbConfig)
	case "accountAdd":
		// name, password
		err = adminAdd(dbConfig, cfg.Args.Num(1), cfg.Args.Num(2))
	default:
		err = errors.New("Must specify a command")
	}

	if err != nil {
		return err
	}

	return nil
}

func migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return err
	}

	fmt.Println("Migrations complete")
	return nil
}

func seed(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Seed(db); err != nil {
		return err
	}

	fmt.Println("Seed data complete")
	return nil
}

func adminAdd(cfg database.Config, name, password string) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if password == "" {
		return errors.New("adminAdd command must be called with additional argument password")
	}

	fmt.Printf("Admin account will be created with password %q\n", password)
	fmt.Print("Continue? (1/0) ")

	var confirm bool
	if _, err := fmt.Scanf("%t\n", &confirm); err != nil {
		return errors.Wrap(err, "processing response")
	}

	if !confirm {
		fmt.Println("Canceling")
		return nil
	}

	ctx := context.Background()

	na := account.NewAccount{
		Name:            name,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleStation},
	}

	a, err := account.Create(ctx, db, na, time.Now())
	if err != nil {
		return err
	}

	fmt.Println("Account created with id:", a.Id)
	return nil
}
