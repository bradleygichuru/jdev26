package main

import (
	"go-rdbms/engine"
	"go-rdbms/parser"
	"testing"
)

func TestBasicCRUD(t *testing.T) {
	db := engine.NewDatabase()

	// Test CREATE TABLE
	stmt := &parser.CreateTableStatement{
		TableName: "users",
		Columns: []*parser.ColumnDefinition{
			{Name: "id", DataType: parser.DATATYPE_INTEGER, PrimaryKey: true},
			{Name: "name", DataType: parser.DATATYPE_TEXT},
		},
	}

	err := db.ExecuteCreateTable(stmt)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Test INSERT
	insertStmt := &parser.InsertStatement{
		TableName: "users",
		Values: []parser.Expression{
			&parser.Literal{Value: 1, Type: parser.DATATYPE_INTEGER},
			&parser.Literal{Value: "Alice", Type: parser.DATATYPE_TEXT},
		},
	}

	err = db.ExecuteInsert(insertStmt)
	if err != nil {
		t.Fatalf("Failed to insert row: %v", err)
	}

	// Test SELECT
	selectStmt := &parser.SelectStatement{
		TableName: "users",
		Columns:   []parser.Expression{&parser.StarExpression{}},
	}

	result, err := db.ExecuteSelect(selectStmt)
	if err != nil {
		t.Fatalf("Failed to select rows: %v", err)
	}

	if len(result.Rows) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(result.Rows))
	}

	if result.Rows[0][0] != 1 || result.Rows[0][1] != "Alice" {
		t.Fatalf("Unexpected row data: %v", result.Rows[0])
	}
}

func TestPrimaryKeyConstraint(t *testing.T) {
	db := engine.NewDatabase()

	// Create table
	stmt := &parser.CreateTableStatement{
		TableName: "test",
		Columns: []*parser.ColumnDefinition{
			{Name: "id", DataType: parser.DATATYPE_INTEGER, PrimaryKey: true},
		},
	}

	db.ExecuteCreateTable(stmt)

	// Insert first row
	insert1 := &parser.InsertStatement{
		TableName: "test",
		Values:    []parser.Expression{&parser.Literal{Value: 1, Type: parser.DATATYPE_INTEGER}},
	}
	err := db.ExecuteInsert(insert1)
	if err != nil {
		t.Fatalf("First insert should succeed: %v", err)
	}

	// Try to insert duplicate primary key
	insert2 := &parser.InsertStatement{
		TableName: "test",
		Values:    []parser.Expression{&parser.Literal{Value: 1, Type: parser.DATATYPE_INTEGER}},
	}
	err = db.ExecuteInsert(insert2)
	if err == nil {
		t.Fatal("Second insert should fail due to primary key violation")
	}
}

func TestParser(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"CREATE TABLE users (id INTEGER PRIMARY KEY)", "CREATE TABLE users (id INTEGER PRIMARY KEY)"},
		{"INSERT INTO users VALUES (1, 'Alice')", "INSERT INTO users VALUES (1, 'Alice')"},
		{"SELECT * FROM users", "SELECT * FROM users"},
	}

	for _, test := range tests {
		lexer := parser.NewLexer(test.input)
		p := parser.NewParser(lexer)
		stmt, err := p.ParseStatement()
		if err != nil {
			t.Fatalf("Parse error for %s: %v", test.input, err)
		}
		if stmt.String() != test.expected {
			t.Fatalf("Expected %s, got %s", test.expected, stmt.String())
		}
	}
}
