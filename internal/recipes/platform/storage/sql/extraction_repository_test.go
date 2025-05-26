package sql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	recipesdomain "github.com/rubenbupe/recipe-video-parser/internal/recipes/domain"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ExtractionRepository_Save_RepositoryError(t *testing.T) {
	extractionID, userID, data, metadata, createdAt := "37a0f027-15e6-47cc-a5d2-64183281087e", "37a0f027-15e6-47cc-a5d2-64183281087e", "{\"field\":\"value\"}", "{\"meta\":\"value\"}", "2023-10-01T00:00:00Z"
	extraction, err := recipesdomain.NewExtraction(extractionID, userID, data, metadata, createdAt)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO recipe_extractions (id, user_id, data, metadata, created_at) VALUES (?, ?, ?, ?, ?)").
		WithArgs(extractionID, userID, data, metadata, createdAt).
		WillReturnError(errors.New("something-failed"))

	repo := NewExtractionRepository(&connection, &config)

	err = repo.Save(context.Background(), extraction)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}

func Test_ExtractionRepository_Save_Succeed(t *testing.T) {
	extractionID, userID, data, metadata, createdAt := "37a0f027-15e6-47cc-a5d2-64183281087e", "37a0f027-15e6-47cc-a5d2-64183281087e", "{\"field\":\"value\"}", "{\"meta\":\"value\"}", "2023-10-01T00:00:00Z"

	extraction, err := recipesdomain.NewExtraction(extractionID, userID, data, metadata, createdAt)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO recipe_extractions (id, user_id, data, metadata, created_at) VALUES (?, ?, ?, ?, ?)").
		WithArgs(extractionID, userID, data, metadata, createdAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewExtractionRepository(&connection, &config)

	err = repo.Save(context.Background(), extraction)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
}

func Test_ExtractionRepository_Exists_RepositoryError(t *testing.T) {
	id := "37a0f027-15e6-47cc-a5d2-64183281087e"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT 1 FROM recipe_extractions WHERE id = ?").
		WithArgs(id).
		WillReturnError(errors.New("something-failed"))

	repo := NewExtractionRepository(&connection, &config)

	extractionID, err := recipesdomain.NewExtractionID(id)
	require.NoError(t, err)
	exists, err := repo.Exists(context.Background(), extractionID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
	assert.False(t, exists)
}

func Test_ExtractionRepository_Exists_Succeed(t *testing.T) {
	id := "37a0f027-15e6-47cc-a5d2-64183281087e"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT 1 FROM recipe_extractions WHERE id = ?").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	repo := NewExtractionRepository(&connection, &config)

	extractionID, err := recipesdomain.NewExtractionID(id)
	require.NoError(t, err)
	exists, err := repo.Exists(context.Background(), extractionID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.True(t, exists)
}

func Test_ExtractionRepository_Get_RepositoryError(t *testing.T) {
	id := "37a0f027-15e6-47cc-a5d2-64183281087e"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT id, user_id, data, metadata, created_at FROM recipe_extractions WHERE id = ?").
		WithArgs(id).
		WillReturnError(errors.New("something-failed"))

	repo := NewExtractionRepository(&connection, &config)

	extractionID, err := recipesdomain.NewExtractionID(id)
	require.NoError(t, err)
	extraction, err := repo.Get(context.Background(), extractionID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
	assert.Nil(t, extraction)
}

func Test_ExtractionRepository_Get_NotFound(t *testing.T) {
	id := "37a0f027-15e6-47cc-a5d2-64183281087e"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT id, user_id, data, metadata, created_at FROM recipe_extractions WHERE id = ?").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "data", "metadata", "created_at"}))

	repo := NewExtractionRepository(&connection, &config)

	extractionID, err := recipesdomain.NewExtractionID(id)
	require.NoError(t, err)
	extraction, err := repo.Get(context.Background(), extractionID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Nil(t, extraction)
}

func Test_ExtractionRepository_Get_Succeed(t *testing.T) {
	id := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	data := "{\"field\":\"value\"}"
	metadata := "{\"meta\":\"value\"}"
	createdAt := "2023-10-01T00:00:00Z"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT id, user_id, data, metadata, created_at FROM recipe_extractions WHERE id = ?").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "data", "metadata", "created_at"}).AddRow(id, userID, data, metadata, createdAt))

	repo := NewExtractionRepository(&connection, &config)

	extractionID, err := recipesdomain.NewExtractionID(id)
	require.NoError(t, err)
	extraction, err := repo.Get(context.Background(), extractionID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.NotNil(t, extraction)
}

func Test_ExtractionRepository_GetByUserID_RepositoryError(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT id, user_id, data, metadata, created_at FROM recipe_extractions WHERE user_id = ?").
		WithArgs(userID).
		WillReturnError(errors.New("something-failed"))

	repo := NewExtractionRepository(&connection, &config)

	extractionUserID, err := recipesdomain.NewExtractionUserID(userID)
	require.NoError(t, err)
	extractions, err := repo.GetByUserID(context.Background(), extractionUserID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
	assert.Nil(t, extractions)
}

func Test_ExtractionRepository_GetByUserID_NotFound(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT id, user_id, data, metadata, created_at FROM recipe_extractions WHERE user_id = ?").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "data", "metadata", "created_at"}))

	repo := NewExtractionRepository(&connection, &config)

	extractionUserID, err := recipesdomain.NewExtractionUserID(userID)
	require.NoError(t, err)
	extractions, err := repo.GetByUserID(context.Background(), extractionUserID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Empty(t, extractions)
}

func Test_ExtractionRepository_GetByUserID_Succeed(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	id := "37a0f027-15e6-47cc-a5d2-64183281087e"
	data := "{\"field\":\"value\"}"
	metadata := "{\"meta\":\"value\"}"
	createdAt := "2023-10-01T00:00:00Z"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT id, user_id, data, metadata, created_at FROM recipe_extractions WHERE user_id = ?").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "data", "metadata", "created_at"}).AddRow(id, userID, data, metadata, createdAt))

	repo := NewExtractionRepository(&connection, &config)

	extractionUserID, err := recipesdomain.NewExtractionUserID(userID)
	require.NoError(t, err)
	extractions, err := repo.GetByUserID(context.Background(), extractionUserID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Len(t, extractions, 1)
	assert.Equal(t, id, extractions[0].Id.String())
	assert.Equal(t, userID, extractions[0].UserId.String())
	assert.Equal(t, data, extractions[0].Data)
	assert.Equal(t, metadata, extractions[0].Metadata)
	assert.Equal(t, createdAt, extractions[0].CreatedAt.String())
}
