package db

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetDSN creates a database connection string from environment variables.
func GetDSN() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "test"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "test"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "test"
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)
}

// CreateDatabaseAndTable connects to the database and ensures the necessary table exists.
func CreateDatabaseAndTable(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := GetDSN()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Printf("Warning: Database might not exist. Attempting to create it.")
		// Attempt to create the database first
		dsnWithoutDB := fmt.Sprintf("host=%s port=%s user=%s password=%s",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
		dbConn, err := pgx.Connect(ctx, dsnWithoutDB)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to 'postgres' to create new one: %w", err)
		}
		defer dbConn.Close(ctx)

		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "test"
		}
		_, err = dbConn.Exec(ctx, "CREATE DATABASE "+dbName)
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
		// Now connect to the newly created database
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to newly created database: %w", err)
		}
	}

	// Ping to check connection
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	// Create the ip_country table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS ip_country (
		ip VARCHAR(45) PRIMARY KEY,
		country VARCHAR(2) NOT NULL
	);`
	_, err = pool.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Database connection and table check successful.")
	return pool, nil
}

// ValidateIP checks if the string is a valid IPv4 or IPv6 address.
func ValidateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil
}

// SaveCountryToDB inserts or updates an IP-to-country mapping.
func SaveCountryToDB(ctx context.Context, db *pgxpool.Pool, ip, country string) error {
	query := `
	INSERT INTO ip_country (ip, country)
	VALUES ($1, $2)
	ON CONFLICT (ip) DO UPDATE SET country = $2;`

	_, err := db.Exec(ctx, query, ip, country)
	if err != nil {
		return fmt.Errorf("failed to save country to db: %w", err)
	}
	return nil
}

// GetCountryFromDB retrieves a country for a given IP from the database.
func GetCountryFromDB(ctx context.Context, db *pgxpool.Pool, ip string) (string, error) {
	var country string
	query := "SELECT country FROM ip_country WHERE ip = $1;"
	err := db.QueryRow(ctx, query, ip).Scan(&country)
	if err != nil {
		return "", fmt.Errorf("failed to get country from db: %w", err)
	}
	return country, nil
}
