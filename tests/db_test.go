package tests

import (
	"sync"

	"github.com/oloomoses/opinions-hub/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var TestDB *gorm.DB
var once sync.Once

func init() {
	once.Do(func() {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

		if err != nil {
			panic("Failed to init test db: " + err.Error())
		}

		err = db.AutoMigrate(
			&models.Opinion{},
			&models.Image{},
			&models.User{},
		)

		if err != nil {
			panic("Failed to migrate test db: " + err.Error())
		}

		TestDB = db
	})

}
