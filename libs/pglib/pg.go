package pglib

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
)

// PG defines the interface for database operations.
type PG interface {
	NewTable(tableName string, schema string) error
	Insert(tableName string, data map[string]interface{}) error
	InsertBulk(tableName string, data []map[string]interface{}) error
	Select(query string) []map[string]interface{}
	Delete(query string) error
	ShowAllTables() ([]string, error)
	ShowAllDBs() ([]string, error)
	Print(data []map[string]interface{})
}

// PGImpl implements the PG interface using pgx.
type PGImpl struct {
	pool *pgxpool.Pool
}

// New creates a new PGImpl instance.
func New(connString string) (PG, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("connecting to PostgreSQL: %w", err)
	}
	return &PGImpl{pool: pool}, nil
}

// NewTable creates a new table with the given schema.
func (db *PGImpl) NewTable(tableName string, schema string) error {
	sqlStatement := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, schema)
	_, err := db.pool.Exec(context.Background(), sqlStatement)
	return err
}

// Insert inserts data into the specified table.
func (db *PGImpl) Insert(tableName string, data map[string]interface{}) error {
	columns := ""
	values := ""
	args := []interface{}{}
	i := 1
	for k, v := range data {
		columns += fmt.Sprintf("%s,", k)
		values += fmt.Sprintf("$%d,", i)
		args = append(args, v)
		i++
	}
	columns = columns[:len(columns)-1]
	values = values[:len(values)-1]
	sqlStatement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, columns, values)
	_, err := db.pool.Exec(context.Background(), sqlStatement, args...)
	return err
}

// Select performs a raw SQL query and returns the results.
func (db *PGImpl) Select(query string) []map[string]interface{} {
	rows, err := db.pool.Query(context.Background(), query)
	if err != nil {
		fmt.Println("postgres: ", err)
		return nil
	}
	defer rows.Close()
	columns := rows.FieldDescriptions()
	var results []map[string]interface{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			fmt.Println("postgres: ", err)
			return nil
		}
		result := make(map[string]interface{})
		for i, v := range values {
			result[string(columns[i].Name)] = v
		}
		results = append(results, result)
	}
	return results
}

// Delete performs a raw SQL delete operation.
func (db *PGImpl) Delete(query string) error {
	_, err := db.pool.Exec(context.Background(), query)
	return err
}

// ShowAllTables lists all tables in the database.
func (db *PGImpl) ShowAllTables() ([]string, error) {
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema='public';"
	rows, err := db.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}

// ShowAllDBs lists all databases.
func (db *PGImpl) ShowAllDBs() ([]string, error) {
	query := "SELECT datname FROM pg_database WHERE datistemplate = false;"
	rows, err := db.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dbs []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		dbs = append(dbs, dbName)
	}
	return dbs, nil
}

func (db *PGImpl) InsertBulk(tableName string, data []map[string]interface{}) error {
	for _, record := range data {
		if err := db.Insert(tableName, record); err != nil {
			return err
		}
	}
	return nil
}

func (db *PGImpl) Print(data []map[string]interface{}) {
	if len(data) == 0 {
		return
	}

	// Get column names
	var columns []string
	for k := range data[0] {
		columns = append(columns, k)
	}

	// Print header
	fmt.Println(strings.Join(columns, ", "))

	// Print data
	for _, row := range data {
		var values []string
		for _, col := range columns {
			value := row[col]
			switch v := value.(type) {
			case string:
				values = append(values, fmt.Sprintf("%q", v))
			default:
				values = append(values, fmt.Sprintf("%v", v))
			}
		}
		fmt.Println(strings.Join(values, ", "))
	}
}
