package database

import (
	"database/sql"
)

// DBInterface allows for mocking database operations
type DBInterface interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// MockDB is a mock implementation for testing
type MockDB struct {
	QueryRowFunc func(query string, args ...interface{}) *sql.Row
	QueryFunc    func(query string, args ...interface{}) (*sql.Rows, error)
	ExecFunc     func(query string, args ...interface{}) (sql.Result, error)
}

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(query, args...)
	}
	return nil
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(query, args...)
	}
	return nil, nil
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(query, args...)
	}
	return nil, nil
}
