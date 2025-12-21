package repository_test

import (
	"testing"

	"github.com/oloomoses/opinions-hub/internal/models"
	"github.com/oloomoses/opinions-hub/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Opinion{})
	return db
}

func TestOpinionRepository_Create(t *testing.T) {
	db := testDB()
	opinions := []struct {
		content string
		wantErr bool
	}{
		{
			content: "This content is not null and should save to db",
			wantErr: false,
		},

		{
			content: "",
			wantErr: true,
		},
	}

	for _, tt := range opinions {
		t.Run(tt.content, func(t *testing.T) {
			tx := db.Begin()
			opinionRepo := repository.NewOpinionRepo(tx)

			defer tx.Rollback()

			op := &models.Opinion{Content: tt.content}

			err := opinionRepo.Create(op)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// var created models.Opinion

			opin, err := opinionRepo.GetByID(uint(op.ID))
			// tx.First(&created, "content = ?", tt.content)

			// assert.Equal(t, tt.content, created.Content)
			assert.NoError(t, err)
			assert.NotZero(t, opin.ID)

		})
	}

}

func TestOpinionRepo_GetAll(t *testing.T) {
	db := testDB()

	tx := db.Begin()
	defer tx.Rollback()
	opnRepo := repository.NewOpinionRepo(tx)

	opinions := []string{
		"First Opinion",
		"Second Opinion",
		"Third Opinion",
		"Fourth Opinion",
		"Fifth Opinion",
		"Sixth Opinion",
		"Seventh Opinion",
	}

	for _, opn := range opinions {
		err := opnRepo.Create(&models.Opinion{Content: opn})

		assert.NoError(t, err)
	}

	opns, err := opnRepo.GetAll()

	assert.NoError(t, err)

	assert.Len(t, opns, len(opinions))

	for i, p := range opinions {
		assert.Equal(t, p, opns[i].Content)
	}

}

func TestOpinionRepo_GetByID(t *testing.T) {
	db := testDB()

	tx := db.Begin()
	opinionRepo := repository.NewOpinionRepo(tx)

	defer tx.Rollback()

	opn := &models.Opinion{Content: "Find me!"}

	err := opinionRepo.Create(opn)

	assert.NoError(t, err)
	assert.NotZero(t, opn.ID)

	found, err := opinionRepo.GetByID(uint(opn.ID))

	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, opn.ID, found.ID)
	assert.Equal(t, opn.Content, found.Content)
}
