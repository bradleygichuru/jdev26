package parser

import (
	"errors"
	"fmt"
	"strconv"
)

// Parser converts tokens into AST nodes
type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []string
}

// NewParser creates a new parser
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

// ParseStatement parses a single SQL statement
func (p *Parser) ParseStatement() (Statement, error) {
	switch p.currentToken.Type {
	case TOKEN_SELECT:
		return p.parseSelectStatement()
	case TOKEN_INSERT:
		return p.parseInsertStatement()
	case TOKEN_UPDATE:
		return p.parseUpdateStatement()
	case TOKEN_DELETE:
		return p.parseDeleteStatement()
	case TOKEN_CREATE:
		return p.parseCreateStatement()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.currentToken.Literal)
	}
}

// parseCreateStatement parses CREATE TABLE statements
func (p *Parser) parseCreateStatement() (*CreateTableStatement, error) {
	stmt := &CreateTableStatement{}

	if !p.expectPeek(TOKEN_TABLE) {
		return nil, errors.New("expected TABLE after CREATE")
	}

	if !p.expectPeek(TOKEN_IDENTIFIER) {
		return nil, errors.New("expected table name after TABLE")
	}
	stmt.TableName = p.currentToken.Literal

	if !p.expectPeek(TOKEN_LEFT_PAREN) {
		return nil, errors.New("expected ( after table name")
	}

	stmt.Columns = p.parseColumnDefinitions()

	if !p.expectPeek(TOKEN_RIGHT_PAREN) {
		return nil, errors.New("expected ) after column definitions")
	}

	return stmt, nil
}

// parseColumnDefinitions parses column definitions in CREATE TABLE
func (p *Parser) parseColumnDefinitions() []*ColumnDefinition {
	var columns []*ColumnDefinition

	for !p.peekTokenIs(TOKEN_RIGHT_PAREN) && !p.peekTokenIs(TOKEN_EOF) {
		col := &ColumnDefinition{}

		if !p.expectPeek(TOKEN_IDENTIFIER) {
			break
		}
		col.Name = p.currentToken.Literal

		dataType, err := p.parseDataType()
		if err != nil {
			break
		}
		col.DataType = dataType

		// Check for PRIMARY KEY or UNIQUE constraints
		if p.peekTokenIs(TOKEN_PRIMARY) {
			p.nextToken()
			if p.expectPeek(TOKEN_KEY) {
				col.PrimaryKey = true
			}
		} else if p.peekTokenIs(TOKEN_UNIQUE) {
			p.nextToken()
			col.Unique = true
		}

		columns = append(columns, col)

		if !p.peekTokenIs(TOKEN_RIGHT_PAREN) {
			if !p.expectPeek(TOKEN_COMMA) {
				break
			}
		}
	}

	return columns
}

// parseDataType parses data type specifications
func (p *Parser) parseDataType() (DataType, error) {
	if !p.expectPeek(TOKEN_IDENTIFIER) {
		return DATATYPE_INTEGER, errors.New("expected data type")
	}

	switch p.currentToken.Literal {
	case "INTEGER", "INT":
		return DATATYPE_INTEGER, nil
	case "TEXT", "VARCHAR":
		return DATATYPE_TEXT, nil
	case "BOOLEAN", "BOOL":
		return DATATYPE_BOOLEAN, nil
	default:
		return DATATYPE_INTEGER, fmt.Errorf("unknown data type: %s", p.currentToken.Literal)
	}
}

// parseInsertStatement parses INSERT statements
func (p *Parser) parseInsertStatement() (*InsertStatement, error) {
	stmt := &InsertStatement{}

	if !p.expectPeek(TOKEN_INTO) {
		return nil, errors.New("expected INTO after INSERT")
	}

	if !p.expectPeek(TOKEN_IDENTIFIER) {
		return nil, errors.New("expected table name after INTO")
	}
	stmt.TableName = p.currentToken.Literal

	if !p.expectPeek(TOKEN_VALUES) {
		return nil, errors.New("expected VALUES after table name")
	}

	if !p.expectPeek(TOKEN_LEFT_PAREN) {
		return nil, errors.New("expected ( after VALUES")
	}

	stmt.Values = p.parseExpressionList(TOKEN_RIGHT_PAREN)

	if !p.expectPeek(TOKEN_RIGHT_PAREN) {
		return nil, errors.New("expected ) after values")
	}

	return stmt, nil
}

// parseSelectStatement parses SELECT statements
func (p *Parser) parseSelectStatement() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	stmt.Columns = p.parseSelectColumns()

	if !p.expectPeek(TOKEN_FROM) {
		return nil, errors.New("expected FROM after SELECT columns")
	}

	if !p.expectPeek(TOKEN_IDENTIFIER) {
		return nil, errors.New("expected table name after FROM")
	}
	stmt.TableName = p.currentToken.Literal

	// Optional JOIN clause
	if p.peekTokenIs(TOKEN_JOIN) {
		join, err := p.parseJoinClause()
		if err != nil {
			return nil, err
		}
		stmt.Join = join
	}

	// Optional WHERE clause
	if p.peekTokenIs(TOKEN_WHERE) {
		p.nextToken()
		where, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	return stmt, nil
}

// parseSelectColumns parses column list in SELECT
func (p *Parser) parseSelectColumns() []Expression {
	var columns []Expression

	if p.peekTokenIs(TOKEN_STAR) {
		p.nextToken()
		columns = append(columns, &StarExpression{})
		return columns
	}

	columns = p.parseExpressionList(TOKEN_FROM)
	return columns
}

// parseJoinClause parses JOIN clause
func (p *Parser) parseJoinClause() (*JoinClause, error) {
	join := &JoinClause{}

	p.nextToken() // consume JOIN

	if !p.expectPeek(TOKEN_IDENTIFIER) {
		return nil, errors.New("expected table name after JOIN")
	}
	join.TableName = p.currentToken.Literal

	if !p.expectPeek(TOKEN_ON) {
		return nil, errors.New("expected ON after JOIN table")
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// For simplicity, assume it's a binary expression
	if binaryExpr, ok := expr.(*BinaryExpression); ok {
		join.On = binaryExpr
	} else {
		return nil, errors.New("expected binary expression in ON clause")
	}

	return join, nil
}

// parseUpdateStatement parses UPDATE statements
func (p *Parser) parseUpdateStatement() (*UpdateStatement, error) {
	stmt := &UpdateStatement{}

	if !p.expectPeek(TOKEN_IDENTIFIER) {
		return nil, errors.New("expected table name after UPDATE")
	}
	stmt.TableName = p.currentToken.Literal

	if !p.expectPeek(TOKEN_SET) {
		return nil, errors.New("expected SET after table name")
	}

	stmt.Set = p.parseSetClause()

	// Optional WHERE clause
	if p.peekTokenIs(TOKEN_WHERE) {
		p.nextToken()
		where, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	return stmt, nil
}

// parseSetClause parses SET column = value pairs
func (p *Parser) parseSetClause() map[string]Expression {
	set := make(map[string]Expression)

	for !p.peekTokenIs(TOKEN_WHERE) && !p.peekTokenIs(TOKEN_EOF) {
		if !p.expectPeek(TOKEN_IDENTIFIER) {
			break
		}
		colName := p.currentToken.Literal

		if !p.expectPeek(TOKEN_EQUALS) {
			break
		}

		expr, err := p.parseExpression()
		if err != nil {
			break
		}

		set[colName] = expr

		if !p.peekTokenIs(TOKEN_WHERE) {
			if !p.expectPeek(TOKEN_COMMA) {
				break
			}
		}
	}

	return set
}

// parseDeleteStatement parses DELETE statements
func (p *Parser) parseDeleteStatement() (*DeleteStatement, error) {
	stmt := &DeleteStatement{}

	if !p.expectPeek(TOKEN_FROM) {
		return nil, errors.New("expected FROM after DELETE")
	}

	if !p.expectPeek(TOKEN_IDENTIFIER) {
		return nil, errors.New("expected table name after FROM")
	}
	stmt.TableName = p.currentToken.Literal

	// Optional WHERE clause
	if p.peekTokenIs(TOKEN_WHERE) {
		p.nextToken()
		where, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		stmt.Where = where
	}

	return stmt, nil
}

// parseExpressionList parses comma-separated expression lists
func (p *Parser) parseExpressionList(endToken TokenType) []Expression {
	var expressions []Expression

	for !p.peekTokenIs(endToken) && !p.peekTokenIs(TOKEN_EOF) {
		expr, err := p.parseExpression()
		if err != nil {
			break
		}
		expressions = append(expressions, expr)

		if !p.peekTokenIs(endToken) {
			if !p.expectPeek(TOKEN_COMMA) {
				break
			}
		}
	}

	return expressions
}

// parseExpression parses expressions (simplified version)
func (p *Parser) parseExpression() (Expression, error) {
	left, err := p.parsePrimaryExpression()
	if err != nil {
		return nil, err
	}

	// Check for binary operators
	if p.peekTokenIs(TOKEN_EQUALS) || p.peekTokenIs(TOKEN_NOT_EQUALS) ||
		p.peekTokenIs(TOKEN_GREATER) || p.peekTokenIs(TOKEN_LESS) ||
		p.peekTokenIs(TOKEN_GREATER_EQUALS) || p.peekTokenIs(TOKEN_LESS_EQUALS) {

		p.nextToken()
		operator := p.currentToken.Literal

		right, err := p.parsePrimaryExpression()
		if err != nil {
			return nil, err
		}

		return &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}, nil
	}

	return left, nil
}

// parsePrimaryExpression parses primary expressions (literals, identifiers)
func (p *Parser) parsePrimaryExpression() (Expression, error) {
	switch p.peekToken.Type {
	case TOKEN_IDENTIFIER:
		p.nextToken()
		ident := &Identifier{Value: p.currentToken.Literal}

		// Check for qualified identifier (table.column)
		if p.peekTokenIs(TOKEN_DOT) {
			p.nextToken() // consume dot
			if !p.expectPeek(TOKEN_IDENTIFIER) {
				return nil, errors.New("expected identifier after dot")
			}
			return &QualifiedIdentifier{
				Table:  ident.Value,
				Column: p.currentToken.Literal,
			}, nil
		}

		return ident, nil
	case TOKEN_STRING:
		p.nextToken()
		return &Literal{Value: p.currentToken.Literal, Type: DATATYPE_TEXT}, nil
	case TOKEN_NUMBER:
		p.nextToken()
		value, err := strconv.Atoi(p.currentToken.Literal)
		if err != nil {
			return nil, err
		}
		return &Literal{Value: value, Type: DATATYPE_INTEGER}, nil
	case TOKEN_TRUE, TOKEN_FALSE:
		p.nextToken()
		value := p.currentToken.Type == TOKEN_TRUE
		return &Literal{Value: value, Type: DATATYPE_BOOLEAN}, nil
	default:
		return nil, fmt.Errorf("unexpected token in expression: %s", p.peekToken.Literal)
	}
}

// Helper methods
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %v, got %v instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// GetErrors returns parser errors
func (p *Parser) GetErrors() []string {
	return p.errors
}
