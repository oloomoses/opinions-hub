package tests

import (
	"os"
	"testing"

	"github.com/oloomoses/opinions-hub/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var GlobalDB *gorm.DB

func TestMain(m *testing.M) {
	var err error

	GlobalDB, err = gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{},
	)

	if err != nil {
		panic(err)
	}

	// Migrate Models once
	GlobalDB.AutoMigrate(
		&models.Opinion{},
	)

	// run tests
	code := m.Run()

	os.Exit(code)
}
