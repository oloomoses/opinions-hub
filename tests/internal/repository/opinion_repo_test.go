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

			opin, err := opinionRepo.GetByID(uint(op.ID))
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

func TestOpinionRepo_Update(t *testing.T) {
	db := testDB()

	tx := db.Begin()

	defer tx.Rollback()

	opinionRepo := repository.NewOpinionRepo(tx)

	savedOpn := &models.Opinion{Content: "Original Opinion"}

	err := opinionRepo.Create(savedOpn)

	assert.NoError(t, err)
	assert.NotZero(t, savedOpn.ID)

	update := make(map[string]interface{})

	update["content"] = "Updated Opinion"

	testCase := []struct {
		tcase   string
		id      uint
		update  map[string]interface{}
		wantErr bool
	}{
		{
			tcase:   "Successfull Update",
			id:      uint(savedOpn.ID),
			update:  update,
			wantErr: false,
		},
		{
			tcase:   "Unsuccesful Update",
			id:      9999,
			update:  update,
			wantErr: true,
		},
	}

	for _, tt := range testCase {
		t.Run(tt.tcase, func(t *testing.T) {
			err := opinionRepo.Update(uint(tt.id), tt.update)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
				return
			}

			assert.NoError(t, err)

			var updatedOpn models.Opinion

			err = tx.First(&updatedOpn, tt.id).Error
			assert.NoError(t, err)
			assert.Equal(t, "Updated Opinion", updatedOpn.Content)

		})
	}
}

func TestOpinionRepo_Delete(t *testing.T) {
	db := testDB()

	delTests := []struct {
		name    string
		setup   func(repo repository.OpinionRepo) uint
		id      uint
		wantErr bool
	}{
		{
			name: "Success if id is exitst",
			setup: func(repo repository.OpinionRepo) uint {
				opn := &models.Opinion{Content: "Opinion to delete"}

				_ = repo.Create(opn)
				return uint(opn.ID)
			},
			wantErr: false,
		},

		{
			name:    "Fail if id not found",
			setup:   nil,
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range delTests {
		t.Run(tt.name, func(t *testing.T) {

			tx := db.Begin()
			defer tx.Rollback()

			repo := repository.NewOpinionRepo(tx)

			id := tt.id

			if tt.setup != nil {
				id = tt.setup(repo)
			}
			err := repo.Delete(id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			var opinion models.Opinion

			err = tx.First(&opinion, id).Error

			assert.Error(t, err)
			assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		})
	}
}
