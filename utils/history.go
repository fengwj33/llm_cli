package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName     = "chat_records.db"
	createTable = `
		CREATE TABLE IF NOT EXISTS conversations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			assistant TEXT NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
)

type History struct {
	db *sql.DB
}

type Record struct {
	ID        int64
	Assistant string
	Role      string
	Content   string
}

// NewHistory initializes the history database
func NewHistory() (*History, error) {
	// Get home directory for database storage
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	// Create .config/llm_cli directory if it doesn't exist
	dbDir := filepath.Join(home, ".config", "llm_cli")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Open database connection
	dbPath := filepath.Join(dbDir, dbName)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create table if not exists
	if _, err := db.Exec(createTable); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &History{db: db}, nil
}

// Close closes the database connection
func (h *History) Close() error {
	return h.db.Close()
}

// Push adds a new record to the history
func (h *History) Push(assistant, role, content string) error {
	query := `
		INSERT INTO conversations (assistant, role, content)
		VALUES (?, ?, ?);
	`
	_, err := h.db.Exec(query, assistant, role, content)
	if err != nil {
		return fmt.Errorf("failed to insert record: %v", err)
	}
	return nil
}

// Fetch retrieves the most recent records for a specific assistant
func (h *History) Fetch(assistant string, limit int) ([]Record, error) {
	query := `
		SELECT id, assistant, role, content
		FROM conversations
		WHERE assistant = ?
		ORDER BY created_at DESC, id DESC
		LIMIT ?;
	`
	rows, err := h.db.Query(query, assistant, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch records: %v", err)
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.Assistant, &r.Role, &r.Content); err != nil {
			return nil, fmt.Errorf("failed to scan record: %v", err)
		}
		records = append(records, r)
	}

	// Reverse the records to get chronological order
	for i := 0; i < len(records)/2; i++ {
		j := len(records) - 1 - i
		records[i], records[j] = records[j], records[i]
	}
	return records, nil
}

// Clear removes all history for a specific assistant
func (h *History) Clear(assistant string) error {
	query := `
		DELETE FROM conversations
		WHERE assistant = ?;
	`
	_, err := h.db.Exec(query, assistant)
	if err != nil {
		return fmt.Errorf("failed to clear history: %v", err)
	}
	return nil
} 