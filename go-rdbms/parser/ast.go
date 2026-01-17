package parser

import (
	"fmt"
	"strings"
)

// Node represents any AST node
type Node interface {
	String() string
}

// Statement represents a SQL statement
type Statement interface {
	Node
	statementNode()
}

// Expression represents a SQL expression
type Expression interface {
	Node
	expressionNode()
}

// CreateTableStatement represents CREATE TABLE statement
type CreateTableStatement struct {
	TableName string
	Columns   []*ColumnDefinition
}

func (c *CreateTableStatement) statementNode() {}
func (c *CreateTableStatement) String() string {
	var cols []string
	for _, col := range c.Columns {
		cols = append(cols, col.String())
	}
	return "CREATE TABLE " + c.TableName + " (" + strings.Join(cols, ", ") + ")"
}

// ColumnDefinition represents a column in CREATE TABLE
type ColumnDefinition struct {
	Name       string
	DataType   DataType
	PrimaryKey bool
	Unique     bool
}

func (c *ColumnDefinition) String() string {
	result := c.Name + " " + c.DataType.String()
	if c.PrimaryKey {
		result += " PRIMARY KEY"
	}
	if c.Unique {
		result += " UNIQUE"
	}
	return result
}

// DataType represents SQL data types
type DataType int

const (
	DATATYPE_INTEGER DataType = iota
	DATATYPE_TEXT
	DATATYPE_BOOLEAN
)

func (d DataType) String() string {
	switch d {
	case DATATYPE_INTEGER:
		return "INTEGER"
	case DATATYPE_TEXT:
		return "TEXT"
	case DATATYPE_BOOLEAN:
		return "BOOLEAN"
	default:
		return "UNKNOWN"
	}
}

// InsertStatement represents INSERT statement
type InsertStatement struct {
	TableName string
	Values    []Expression
}

func (i *InsertStatement) statementNode() {}
func (i *InsertStatement) String() string {
	var vals []string
	for _, val := range i.Values {
		vals = append(vals, val.String())
	}
	return "INSERT INTO " + i.TableName + " VALUES (" + strings.Join(vals, ", ") + ")"
}

// SelectStatement represents SELECT statement
type SelectStatement struct {
	TableName string
	Columns   []Expression
	Where     Expression
	Join      *JoinClause
}

func (s *SelectStatement) statementNode() {}
func (s *SelectStatement) String() string {
	var cols []string
	for _, col := range s.Columns {
		cols = append(cols, col.String())
	}

	result := "SELECT " + strings.Join(cols, ", ") + " FROM " + s.TableName
	if s.Join != nil {
		result += " " + s.Join.String()
	}
	if s.Where != nil {
		result += " WHERE " + s.Where.String()
	}
	return result
}

// JoinClause represents JOIN clause
type JoinClause struct {
	TableName string
	On        *BinaryExpression
}

func (j *JoinClause) String() string {
	return "JOIN " + j.TableName + " ON " + j.On.String()
}

// UpdateStatement represents UPDATE statement
type UpdateStatement struct {
	TableName string
	Set       map[string]Expression
	Where     Expression
}

func (u *UpdateStatement) statementNode() {}
func (u *UpdateStatement) String() string {
	var sets []string
	for col, val := range u.Set {
		sets = append(sets, col+" = "+val.String())
	}
	result := "UPDATE " + u.TableName + " SET " + strings.Join(sets, ", ")
	if u.Where != nil {
		result += " WHERE " + u.Where.String()
	}
	return result
}

// DeleteStatement represents DELETE statement
type DeleteStatement struct {
	TableName string
	Where     Expression
}

func (d *DeleteStatement) statementNode() {}
func (d *DeleteStatement) String() string {
	result := "DELETE FROM " + d.TableName
	if d.Where != nil {
		result += " WHERE " + d.Where.String()
	}
	return result
}

// Expressions

// Identifier represents a column or table name
type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

// QualifiedIdentifier represents table.column references
type QualifiedIdentifier struct {
	Table  string
	Column string
}

func (q *QualifiedIdentifier) expressionNode() {}
func (q *QualifiedIdentifier) String() string  { return q.Table + "." + q.Column }

// Literal represents literal values
type Literal struct {
	Value interface{}
	Type  DataType
}

func (l *Literal) expressionNode() {}
func (l *Literal) String() string {
	switch l.Type {
	case DATATYPE_TEXT:
		return "'" + l.Value.(string) + "'"
	case DATATYPE_BOOLEAN:
		if l.Value.(bool) {
			return "TRUE"
		}
		return "FALSE"
	default:
		return fmt.Sprintf("%v", l.Value)
	}
}

// BinaryExpression represents binary operations
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpression) expressionNode() {}
func (b *BinaryExpression) String() string {
	return b.Left.String() + " " + b.Operator + " " + b.Right.String()
}

// StarExpression represents SELECT *
type StarExpression struct{}

func (s *StarExpression) expressionNode() {}
func (s *StarExpression) String() string  { return "*" }
