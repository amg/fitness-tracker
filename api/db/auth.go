package db

import (
	"database/sql"
)

// AuthService represents a service for managing OAuth authentication.
type AuthService struct {
	db *sql.DB
}
