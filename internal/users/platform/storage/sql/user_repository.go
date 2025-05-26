package sql

import (
	"context"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/storage"
	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
)

type UserRepository struct {
	connection *storage.Connection
	dbconfig   *storage.Dbconfig
}

func NewUserRepository(connection *storage.Connection, dbconfig *storage.Dbconfig) *UserRepository {
	return &UserRepository{
		connection: connection,
		dbconfig:   dbconfig,
	}
}

func (r *UserRepository) Save(ctx context.Context, user usersdomain.User) error {
	query := "INSERT INTO " + sqlUserTable + " (id, name, api_key, created_at) VALUES (?, ?, ?, ?) " +
		"ON CONFLICT(id) DO UPDATE SET name=excluded.name, api_key=excluded.api_key, created_at=excluded.created_at"
	args := []interface{}{user.Id.String(), user.Name.String(), user.ApiKey.String(), user.CreatedAt.String()}

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbconfig.Timeout)
	defer cancel()

	_, err := r.connection.Db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to upsert user on database: %v", err)
	}

	return nil
}

func (r *UserRepository) Exists(ctx context.Context, id usersdomain.UserID) (bool, error) {
	sb := sqlbuilder.Select("1").From(sqlUserTable)
	sb.Where(sb.Equal("id", id.String()))
	sb.SetFlavor(sqlbuilder.SQLite)
	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbconfig.Timeout)
	defer cancel()

	rows, err := r.connection.Db.QueryContext(ctxTimeout, query, args...)
	if err != nil {
		return false, fmt.Errorf("error trying to check if user exists on database: %v", err)
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (r *UserRepository) ExistsByName(ctx context.Context, name usersdomain.UserName) (bool, error) {
	sb := sqlbuilder.Select("1").From(sqlUserTable)
	sb.Where(sb.Equal("name", name.String()))
	sb.SetFlavor(sqlbuilder.SQLite)
	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbconfig.Timeout)
	defer cancel()

	rows, err := r.connection.Db.QueryContext(ctxTimeout, query, args...)
	if err != nil {
		return false, fmt.Errorf("error trying to check if user exists by name on database: %v", err)
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (r *UserRepository) Get(ctx context.Context, id usersdomain.UserID) (*usersdomain.User, error) {
	userSQLStruct := sqlbuilder.NewStruct(new(sqlUser))
	sb := sqlbuilder.Select("id", "name", "api_key", "created_at").From(sqlUserTable)
	sb.Where(sb.Equal("id", id.String()))
	sb.SetFlavor(sqlbuilder.SQLite)
	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbconfig.Timeout)
	defer cancel()

	row := r.connection.Db.QueryRowContext(ctxTimeout, query, args...)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("error trying to get user from database: %v", err)
	}
	user := new(sqlUser)
	err := row.Scan(userSQLStruct.Addr(user)...)
	if err != nil {
		return nil, nil
	}

	userVO, err := usersdomain.NewUser(user.ID, user.Name, user.ApiKey, user.CreatedAt)
	return &userVO, err
}

func (r *UserRepository) GetByName(ctx context.Context, name usersdomain.UserName) (*usersdomain.User, error) {
	userSQLStruct := sqlbuilder.NewStruct(new(sqlUser))
	sb := sqlbuilder.Select("id", "name", "api_key", "created_at").From(sqlUserTable)
	sb.Where(sb.Equal("name", name.String()))
	sb.SetFlavor(sqlbuilder.SQLite)
	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbconfig.Timeout)
	defer cancel()

	row := r.connection.Db.QueryRowContext(ctxTimeout, query, args...)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("error trying to get user by name from database: %v", err)
	}
	user := new(sqlUser)
	err := row.Scan(userSQLStruct.Addr(user)...)
	if err != nil {
		return nil, nil
	}

	userVO, err := usersdomain.NewUser(user.ID, user.Name, user.ApiKey, user.CreatedAt)
	return &userVO, err
}

func (r *UserRepository) GetByApiKey(ctx context.Context, apiKey usersdomain.UserApiKey) (*usersdomain.User, error) {
	userSQLStruct := sqlbuilder.NewStruct(new(sqlUser))
	sb := sqlbuilder.Select("id", "name", "api_key", "created_at").From(sqlUserTable)
	sb.Where(sb.Equal("api_key", apiKey.String()))
	sb.SetFlavor(sqlbuilder.SQLite)
	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbconfig.Timeout)
	defer cancel()

	row := r.connection.Db.QueryRowContext(ctxTimeout, query, args...)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("error trying to get user by apiKey from database: %v", err)
	}
	user := new(sqlUser)
	err := row.Scan(userSQLStruct.Addr(user)...)
	if err != nil {
		return nil, nil
	}

	userVO, err := usersdomain.NewUser(user.ID, user.Name, user.ApiKey, user.CreatedAt)
	return &userVO, err
}
