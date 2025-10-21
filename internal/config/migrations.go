package config

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// RunMigrations runs all pending migrations
// This is a simple alternative to golang-migrate tool
// For production, use golang-migrate/migrate instead
func RunMigrations(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")

	// Create users table
	if !db.Migrator().HasTable("users") {
		if err := db.Exec(`
			CREATE TABLE users (
				id BIGSERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				email VARCHAR(255) NOT NULL UNIQUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`).Error; err != nil {
			return fmt.Errorf("failed to create users table: %w", err)
		}
		log.Println("✅ Created users table")
	}

	// Create index on email
	if !db.Migrator().HasIndex("users", "email") {
		if err := db.Migrator().CreateIndex("users", "email"); err != nil {
			return fmt.Errorf("failed to create email index: %w", err)
		}
		log.Println("✅ Created index on users.email")
	}

	log.Println("✅ All migrations completed successfully")
	return nil
}

// RollbackMigrations rolls back all migrations
// Use with caution! This will delete all data
func RollbackMigrations(db *gorm.DB) error {
	log.Println("⚠️  Rolling back migrations...")

	if err := db.Migrator().DropTable("users"); err != nil {
		return fmt.Errorf("failed to drop users table: %w", err)
	}

	log.Println("✅ Migrations rolled back")
	return nil
}

// CheckMigrationStatus returns current migration status
func CheckMigrationStatus(db *gorm.DB) map[string]bool {
	status := make(map[string]bool)

	status["users_table"] = db.Migrator().HasTable("users")
	status["users_email_index"] = db.Migrator().HasIndex("users", "email")

	return status
}