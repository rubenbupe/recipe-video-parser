package sql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/storage"
	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserRepository_Save_RepositoryError(t *testing.T) {
	userID, userName, userApiKey, userCreatedAt := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test User", "test-api-key", "2023-10-01T00:00:00Z"
	user, err := usersdomain.NewUser(userID, userName, userApiKey, userCreatedAt)
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
		"INSERT INTO users (id, name, api_key, created_at) VALUES (?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET name=excluded.name, api_key=excluded.api_key, created_at=excluded.created_at").
		WithArgs(userID, userName, userApiKey, userCreatedAt).
		WillReturnError(errors.New("something-failed"))

	repo := NewUserRepository(&connection, &config)

	err = repo.Save(context.Background(), user)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}

func Test_UserRepository_Save_Succeed(t *testing.T) {
	userID, userName, userApiKey, userCreatedAt := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test User", "test-api-key", "2023-10-01T00:00:00Z"

	user, err := usersdomain.NewUser(userID, userName, userApiKey, userCreatedAt)
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
		"INSERT INTO users (id, name, api_key, created_at) VALUES (?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET name=excluded.name, api_key=excluded.api_key, created_at=excluded.created_at").
		WithArgs(userID, userName, userApiKey, userCreatedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewUserRepository(&connection, &config)

	err = repo.Save(context.Background(), user)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
}

func Test_UserRepository_Exists_RepositoryError(t *testing.T) {
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
		"SELECT 1 FROM users WHERE id = ?").
		WithArgs(id).
		WillReturnError(errors.New("something-failed"))

	repo := NewUserRepository(&connection, &config)

	userID, err := usersdomain.NewUserID(id)
	require.NoError(t, err)
	exists, err := repo.Exists(context.Background(), userID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
	assert.False(t, exists)
}

func Test_UserRepository_Exists_Succeed(t *testing.T) {
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
		"SELECT 1 FROM users WHERE id = ?").
		WithArgs(id).
		WillReturnRows(sqlMock.NewRows([]string{"1"}).AddRow(1))

	repo := NewUserRepository(&connection, &config)

	userID, err := usersdomain.NewUserID(id)
	require.NoError(t, err)
	exists, err := repo.Exists(context.Background(), userID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.True(t, exists)
}

func Test_UserRepository_Get_RepositoryError(t *testing.T) {
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
		"SELECT id, name, api_key, created_at FROM users WHERE id = ?").
		WithArgs(id).
		WillReturnError(errors.New("something-failed"))

	repo := NewUserRepository(&connection, &config)

	userID, err := usersdomain.NewUserID(id)
	require.NoError(t, err)
	user, err := repo.Get(context.Background(), userID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
	assert.Nil(t, user)
}

func Test_UserRepository_Get_NotFound(t *testing.T) {
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
		"SELECT id, name, api_key, created_at FROM users WHERE id = ?").
		WithArgs(id).
		WillReturnRows(sqlMock.NewRows([]string{"id", "name", "api_key", "created_at"}))

	repo := NewUserRepository(&connection, &config)

	userID, err := usersdomain.NewUserID(id)
	require.NoError(t, err)
	user, err := repo.Get(context.Background(), userID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func Test_UserRepository_Get_Succeed(t *testing.T) {
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
		"SELECT id, name, api_key, created_at FROM users WHERE id = ?").
		WithArgs(id).
		WillReturnRows(sqlMock.NewRows([]string{"id", "name", "api_key", "created_at"}).AddRow(id, "Test User", "test-api-key", "2023-10-01T00:00:00Z"))

	repo := NewUserRepository(&connection, &config)

	userID, err := usersdomain.NewUserID(id)
	require.NoError(t, err)
	user, err := repo.Get(context.Background(), userID)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func Test_UserRepository_GetByName_Succeed(t *testing.T) {
	userID, userName, userApiKey, userCreatedAt := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test User", "test-api-key", "2023-10-01T00:00:00Z"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT id, name, api_key, created_at FROM users WHERE name = ?").
		WithArgs(userName).
		WillReturnRows(sqlMock.NewRows([]string{"id", "name", "api_key", "created_at"}).AddRow(userID, userName, userApiKey, userCreatedAt))

	repo := NewUserRepository(&connection, &config)
	userNameVO, err := usersdomain.NewUserName(userName)
	require.NoError(t, err)
	result, err := repo.GetByName(context.Background(), userNameVO)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userName, result.Name.String())
}

func Test_UserRepository_GetByApiKey_Succeed(t *testing.T) {
	userID, userName, userApiKey, userCreatedAt := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test User", "test-api-key", "2023-10-01T00:00:00Z"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT id, name, api_key, created_at FROM users WHERE api_key = ?").
		WithArgs(userApiKey).
		WillReturnRows(sqlMock.NewRows([]string{"id", "name", "api_key", "created_at"}).AddRow(userID, userName, userApiKey, userCreatedAt))

	repo := NewUserRepository(&connection, &config)
	apiKeyVO, err := usersdomain.NewUserApiKey(userApiKey)
	require.NoError(t, err)
	result, err := repo.GetByApiKey(context.Background(), apiKeyVO)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userApiKey, result.ApiKey.String())
}

func Test_UserRepository_ExistsByName_Succeed(t *testing.T) {
	userName := "Test User"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	connection := storage.Connection{
		Db: db,
	}
	config := storage.Dbconfig{
		Timeout: 1 * time.Millisecond,
	}

	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT 1 FROM users WHERE name = ?").
		WithArgs(userName).
		WillReturnRows(sqlMock.NewRows([]string{"1"}).AddRow(1))

	repo := NewUserRepository(&connection, &config)
	userNameVO, err := usersdomain.NewUserName(userName)
	require.NoError(t, err)
	exists, err := repo.ExistsByName(context.Background(), userNameVO)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.True(t, exists)
}
