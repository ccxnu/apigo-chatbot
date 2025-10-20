package dal

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbResult struct {
	Success bool   `db:"success"`
	Code    string `db:"code"`
}

type DAL struct {
	DB *pgxpool.Pool
}

func NewDAL(db *pgxpool.Pool) *DAL {
	return &DAL{DB: db}
}

// QueryRows executes a function that returns multiple rows and scans them into a slice of T.
func QueryRows[T any](d *DAL, ctx context.Context, fnName string, args ...any) ([]T, error) {
	placeholders := generatePlaceholders(len(args))

	query := fmt.Sprintf("SELECT * FROM %s(%s)", pgx.Identifier{fnName}.Sanitize(), placeholders)

	rows, err := d.DB.Query(ctx, query, args...)

	if err != nil {
		return nil, fmt.Errorf("dal.QueryRows: failed to execute query for %s: %w", fnName, err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

// QueryRow executes a function that returns a single row and scans it into an instance of T.
func QueryRow[T any](d *DAL, ctx context.Context, fnName string, args ...any) (*T, error) {
	placeholders := generatePlaceholders(len(args))
	query := fmt.Sprintf("SELECT * FROM %s(%s)", pgx.Identifier{fnName}.Sanitize(), placeholders)

	rows, err := d.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("dal.QueryRow: failed to execute query for %s: %w", fnName, err)
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("dal.QueryRow: failed to scan row for %s: %w", fnName, err)
	}
	return &result, nil
}

// ExecProc executes a stored procedure and scans its output parameters into an instance of T.
// It uses reflection to count ALL fields in T, assuming these correspond to the
// total number of OUT parameters in the PostgreSQL procedure, and that the IN
// parameters follow them.
func ExecProc[T any](d *DAL, ctx context.Context, procName string, args ...any) (*T, error) {
	var result T
	resultType := reflect.TypeOf(result)
	numOutParams := 0

	if resultType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("dal.ExecProc: Result type %T must be a struct", result)
	}

	// Count all fields recursively (to include embedded BaseResult fields)
	for i := 0; i < resultType.NumField(); i++ {
		field := resultType.Field(i)
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			numOutParams += field.Type.NumField()
		} else {
			numOutParams++
		}
	}

	if numOutParams == 0 {
		return nil, fmt.Errorf("dal.ExecProc: Result type %T has no fields to map OUT parameters", result)
	}

	// Create 'NULL' placeholders for every OUT parameter required by the struct T
	outPlaceholders := make([]string, numOutParams)
	for i := range numOutParams {
		outPlaceholders[i] = "NULL"
	}
	outPlaceholderStr := strings.Join(outPlaceholders, ", ")

	// Create positional placeholders ($1, $2, ...) for the IN parameters
	inPlaceholders := generatePlaceholders(len(args))

	// Combine them: (NULL, NULL, NULL, $1, $2, ...)
	var allPlaceholders string
	if len(args) > 0 {
		allPlaceholders = outPlaceholderStr + ", " + inPlaceholders
	} else {
		allPlaceholders = outPlaceholderStr
	}

	query := fmt.Sprintf("CALL %s(%s)", pgx.Identifier{procName}.Sanitize(), allPlaceholders)

	rows, err := d.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("dal.ExecProc: failed to execute procedure %s: %w", procName, err)
	}

	resultScan, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, fmt.Errorf("dal.ExecProc: failed to scan procedure output for %s: %w", procName, err)
	}

	return &resultScan, nil
}

// generatePlaceholders creates a string like "$1, $2, $3".
func generatePlaceholders(count int) string {
	if count == 0 {
		return ""
	}

	placeholders := make([]string, count)

	for i := range count {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	return strings.Join(placeholders, ", ")
}
