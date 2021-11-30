// Package cmd _
package cmd

import (
	"fmt"
	"github.com/jinzhu/gorm"

	"github.com/pressly/goose"
)

var (
	migrationsDir = "/migrations/"
)

// Migrate runs migration from specified folder
func Migrate(db *gorm.DB) error {
	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("error applying mysql driver for migrations: %w", err)
	}

	return goose.Up(db.DB(), migrationsDir)
}
