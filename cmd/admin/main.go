// This program performs administrative tasks for the garage sale service.
package main

import (
	// Core Packages
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"
	"fmt"
	"log" // https://golang.org/pkg/log/
	"os"

	// Third-party packages
	"github.com/pkg/errors"

	// Internal applcation packages
	"github.com/deezone/HydroBytes-BaseStation/internal/account"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/auth"
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
	case "adminAdd":
		// name, password
		err = adminAdd(dbConfig, cfg.Args.Num(1), cfg.Args.Num(2))
	case "keygen":
		err = keygen(cfg.Args.Num(1))
	case "migrate":
		err = migrate(dbConfig)
	case "seed":
		err = seed(dbConfig)
	default:
		err = errors.New("Must specify a command")
	}

	if err != nil {
		return err
	}

	return nil
}

// adminAccount creates an admin account
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

// keygen creates an x509 private key for signing auth tokens.
func keygen(path string) error {
	if path == "" {
		return errors.New("keygen missing argument for key path")
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return errors.Wrap(err, "generating keys")
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "creating private file")
	}
	defer file.Close()

	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	if err := pem.Encode(file, &block); err != nil {
		return errors.Wrap(err, "encoding to private file")
	}

	if err := file.Close(); err != nil {
		return errors.Wrap(err, "closing private file")
	}

	return nil
}

// migrate applies database migrations
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
