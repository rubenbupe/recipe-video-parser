package sql

import (
	"context"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	recipesdomain "github.com/rubenbupe/recipe-video-parser/internal/recipes/domain"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/storage"
)

type ExtractionRepository struct {
	connection *storage.Connection
	dbconfig   *storage.Dbconfig
}

func NewExtractionRepository(connection *storage.Connection, dbconfig *storage.Dbconfig) *ExtractionRepository {
	return &ExtractionRepository{
		connection: connection,
		dbconfig:   dbconfig,
	}
}

func (r *ExtractionRepository) Save(ctx context.Context, extraction recipesdomain.Extraction) error {
	extractionSQLStruct := sqlbuilder.NewStruct(new(sqlExtraction)).For(sqlbuilder.SQLite)
	query, args := extractionSQLStruct.InsertInto(sqlExtractionTable, sqlExtraction{
		ID:        extraction.Id.String(),
		UserID:    extraction.UserId.String(),
		Data:      extraction.Data,
		Metadata:  extraction.Metadata,
		CreatedAt: extraction.CreatedAt.String(),
	}).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbconfig.Timeout)
	defer cancel()

	_, err := r.connection.Db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist extraction on database: %v", err)
	}

	return nil
}

func (r *ExtractionRepository) Exists(ctx context.Context, id recipesdomain.ExtractionID) (bool, error) {
	sb := sqlbuilder.Select("1").From(sqlExtractionTable)
	sb.Where(sb.Equal("id", id.String()))
	sb.SetFlavor(sqlbuilder.SQLite)
	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbconfig.Timeout)
	defer cancel()

	rows, err := r.connection.Db.QueryContext(ctxTimeout, query, args...)
	if err != nil {
		return false, fmt.Errorf("error trying to check if extraction exists on database: %v", err)
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (r *ExtractionRepository) Get(ctx context.Context, id recipesdomain.ExtractionID) (*recipesdomain.Extraction, error) {
	extractionSQLStruct := sqlbuilder.NewStruct(new(sqlExtraction))
	sb := sqlbuilder.Select("id", "user_id", "data", "metadata", "created_at").From(sqlExtractionTable)
	sb.Where(sb.Equal("id", id.String()))
	sb.SetFlavor(sqlbuilder.SQLite)
	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbconfig.Timeout)
	defer cancel()

	row := r.connection.Db.QueryRowContext(ctxTimeout, query, args...)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("error trying to get user from database: %v", err)
	}
	extraction := new(sqlExtraction)
	err := row.Scan(extractionSQLStruct.Addr(extraction)...)
	if err != nil {
		return nil, nil
	}

	extractionVO, err := recipesdomain.NewExtraction(extraction.ID, extraction.UserID, extraction.Data, extraction.Metadata, extraction.CreatedAt)
	return &extractionVO, err
}

func (r *ExtractionRepository) GetByUserID(ctx context.Context, userId recipesdomain.ExtractionUserID) ([]recipesdomain.Extraction, error) {
	sb := sqlbuilder.Select("id", "user_id", "data", "metadata", "created_at").From(sqlExtractionTable)
	sb.Where(sb.Equal("user_id", userId.String()))
	sb.SetFlavor(sqlbuilder.SQLite)
	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbconfig.Timeout)
	defer cancel()

	rows, err := r.connection.Db.QueryContext(ctxTimeout, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error trying to get extractions by user id from database: %v", err)
	}
	defer rows.Close()

	var extractions []recipesdomain.Extraction
	for rows.Next() {
		extractionSQLStruct := sqlbuilder.NewStruct(new(sqlExtraction))
		extraction := new(sqlExtraction)
		if err := rows.Scan(extractionSQLStruct.Addr(extraction)...); err != nil {
			return nil, fmt.Errorf("error scanning extraction row: %v", err)
		}

		extractionVO, err := recipesdomain.NewExtraction(extraction.ID, extraction.UserID, extraction.Data, extraction.Metadata, extraction.CreatedAt)
		if err != nil {
			return nil, err
		}
		extractions = append(extractions, extractionVO)
	}

	return extractions, nil
}
