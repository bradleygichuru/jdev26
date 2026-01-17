package engine

import (
	"fmt"
	"go-rdbms/parser"
	"reflect"
)

// ExecuteCreateTable executes a CREATE TABLE statement
func (db *Database) ExecuteCreateTable(stmt *parser.CreateTableStatement) error {
	if _, exists := db.Tables[stmt.TableName]; exists {
		return fmt.Errorf("table %s already exists", stmt.TableName)
	}

	// Convert parser columns to engine columns
	var columns []*Column
	for _, colDef := range stmt.Columns {
		col := &Column{
			Name:       colDef.Name,
			DataType:   colDef.DataType,
			PrimaryKey: colDef.PrimaryKey,
			Unique:     colDef.Unique,
		}
		columns = append(columns, col)
	}

	table := NewTable(stmt.TableName, columns)
	db.Tables[stmt.TableName] = table

	return nil
}

// ExecuteInsert executes an INSERT statement
func (db *Database) ExecuteInsert(stmt *parser.InsertStatement) error {
	table, exists := db.Tables[stmt.TableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", stmt.TableName)
	}

	if len(stmt.Values) != len(table.Columns) {
		return fmt.Errorf("expected %d values, got %d", len(table.Columns), len(stmt.Values))
	}

	row := NewRow()
	for i, expr := range stmt.Values {
		col := table.Columns[i]
		value, err := db.evaluateExpression(expr)
		if err != nil {
			return err
		}
		row.SetValue(col.Name, value)
	}

	return table.InsertRow(row)
}

// ExecuteSelect executes a SELECT statement
func (db *Database) ExecuteSelect(stmt *parser.SelectStatement) (*ResultSet, error) {
	table, exists := db.Tables[stmt.TableName]
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", stmt.TableName)
	}

	// Build where condition function
	whereCondition := func(row *Row) bool { return true } // default: all rows
	if stmt.Where != nil {
		cond, err := db.buildWhereCondition(stmt.Where)
		if err != nil {
			return nil, err
		}
		whereCondition = cond
	}

	// Handle JOIN if present
	if stmt.Join != nil {
		return db.executeJoinSelect(table, stmt)
	}

	// Get matching rows
	rows := table.SelectRows(nil, whereCondition)

	// Determine columns to return
	var columnNames []string
	if len(stmt.Columns) == 1 {
		if _, ok := stmt.Columns[0].(*parser.StarExpression); ok {
			if stmt.Join != nil {
				// For JOIN queries, include columns from both tables
				rightTable := db.Tables[stmt.Join.TableName]
				columnNames = append(table.GetColumnNames(), rightTable.GetColumnNames()...)
			} else {
				columnNames = table.GetColumnNames()
			}
		} else {
			// Single column
			colName, err := db.extractColumnName(stmt.Columns[0])
			if err != nil {
				return nil, err
			}
			columnNames = []string{colName}
		}
	} else {
		// Multiple columns
		for _, col := range stmt.Columns {
			colName, err := db.extractColumnName(col)
			if err != nil {
				return nil, err
			}
			columnNames = append(columnNames, colName)
		}
	}

	// Build result set
	resultSet := &ResultSet{
		Columns: columnNames,
		Rows:    make([][]interface{}, 0, len(rows)),
	}

	for _, row := range rows {
		var values []interface{}
		for _, colName := range columnNames {
			values = append(values, row.GetValue(colName))
		}
		resultSet.Rows = append(resultSet.Rows, values)
	}

	return resultSet, nil
}

// ExecuteUpdate executes an UPDATE statement
func (db *Database) ExecuteUpdate(stmt *parser.UpdateStatement) error {
	table, exists := db.Tables[stmt.TableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", stmt.TableName)
	}

	// Build updates map
	updates := make(map[string]interface{})
	for colName, expr := range stmt.Set {
		value, err := db.evaluateExpression(expr)
		if err != nil {
			return err
		}
		updates[colName] = value
	}

	// Find rows to update
	var rowsToUpdate []*Row
	whereCondition := func(row *Row) bool { return true }
	if stmt.Where != nil {
		cond, err := db.buildWhereCondition(stmt.Where)
		if err != nil {
			return err
		}
		whereCondition = cond
	}

	for _, row := range table.Rows {
		if whereCondition(row) {
			rowsToUpdate = append(rowsToUpdate, row)
		}
	}

	// Apply updates
	for _, row := range rowsToUpdate {
		if table.PrimaryKey != "" {
			pkValue := row.GetValue(table.PrimaryKey)
			if err := table.UpdateRow(pkValue, updates); err != nil {
				return err
			}
		} else {
			// No primary key - update in place
			for colName, value := range updates {
				row.SetValue(colName, value)
			}
		}
	}

	return nil
}

// ExecuteDelete executes a DELETE statement
func (db *Database) ExecuteDelete(stmt *parser.DeleteStatement) error {
	table, exists := db.Tables[stmt.TableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", stmt.TableName)
	}

	// Find rows to delete
	var rowsToDelete []interface{}
	whereCondition := func(row *Row) bool { return true }
	if stmt.Where != nil {
		cond, err := db.buildWhereCondition(stmt.Where)
		if err != nil {
			return err
		}
		whereCondition = cond
	}

	for _, row := range table.Rows {
		if whereCondition(row) {
			if table.PrimaryKey != "" {
				pkValue := row.GetValue(table.PrimaryKey)
				rowsToDelete = append(rowsToDelete, pkValue)
			}
		}
	}

	// Delete rows
	for _, pkValue := range rowsToDelete {
		if err := table.DeleteRow(pkValue); err != nil {
			return err
		}
	}

	return nil
}

// executeJoinSelect handles SELECT with JOIN
func (db *Database) executeJoinSelect(leftTable *Table, stmt *parser.SelectStatement) (*ResultSet, error) {
	rightTable, exists := db.Tables[stmt.Join.TableName]
	if !exists {
		return nil, fmt.Errorf("joined table %s does not exist", stmt.Join.TableName)
	}

	// Simple nested loop join implementation
	var resultRows [][]interface{}

	leftCol, rightCol, err := db.parseJoinCondition(stmt.Join.On)
	if err != nil {
		return nil, err
	}

	// Get column names for result
	columnNames := append(leftTable.GetColumnNames(), rightTable.GetColumnNames()...)

	for _, leftRow := range leftTable.Rows {
		for _, rightRow := range rightTable.Rows {
			leftValue := leftRow.GetValue(leftCol)
			rightValue := rightRow.GetValue(rightCol)

			if reflect.DeepEqual(leftValue, rightValue) {
				var values []interface{}
				// Add left table columns
				for _, colName := range leftTable.GetColumnNames() {
					values = append(values, leftRow.GetValue(colName))
				}
				// Add right table columns
				for _, colName := range rightTable.GetColumnNames() {
					values = append(values, rightRow.GetValue(colName))
				}
				resultRows = append(resultRows, values)
			}
		}
	}

	return &ResultSet{
		Columns: columnNames,
		Rows:    resultRows,
	}, nil
}

// parseJoinCondition extracts column names from JOIN ON condition
func (db *Database) parseJoinCondition(expr *parser.BinaryExpression) (leftCol, rightCol string, err error) {
	if expr.Operator != "=" {
		return "", "", fmt.Errorf("only equality joins are supported")
	}

	// Extract column name from left side
	leftCol, err = db.extractColumnName(expr.Left)
	if err != nil {
		return "", "", err
	}

	// Extract column name from right side
	rightCol, err = db.extractColumnName(expr.Right)
	if err != nil {
		return "", "", err
	}

	return leftCol, rightCol, nil
}

// extractColumnName extracts column name from identifier or qualified identifier
func (db *Database) extractColumnName(expr parser.Expression) (string, error) {
	switch e := expr.(type) {
	case *parser.Identifier:
		return e.Value, nil
	case *parser.QualifiedIdentifier:
		return e.Column, nil
	default:
		return "", fmt.Errorf("expression must be an identifier or qualified identifier")
	}
}

// buildWhereCondition converts a WHERE expression to a function
func (db *Database) buildWhereCondition(expr parser.Expression) (func(*Row) bool, error) {
	switch e := expr.(type) {
	case *parser.BinaryExpression:
		return db.buildBinaryCondition(e)
	default:
		return nil, fmt.Errorf("unsupported WHERE expression type: %T", expr)
	}
}

// buildBinaryCondition builds a condition function from binary expression
func (db *Database) buildBinaryCondition(expr *parser.BinaryExpression) (func(*Row) bool, error) {
	leftCol, err := db.extractColumnName(expr.Left)
	if err != nil {
		return nil, err
	}

	rightValue, err := db.evaluateExpression(expr.Right)
	if err != nil {
		return nil, err
	}

	return func(row *Row) bool {
		leftValue := row.GetValue(leftCol)
		return db.compareValues(leftValue, rightValue, expr.Operator)
	}, nil
}

// compareValues compares two values using the given operator
func (db *Database) compareValues(left, right interface{}, operator string) bool {
	switch operator {
	case "=":
		return reflect.DeepEqual(left, right)
	case "!=":
		return !reflect.DeepEqual(left, right)
	case ">":
		return compareOrdered(left, right) > 0
	case "<":
		return compareOrdered(left, right) < 0
	case ">=":
		return compareOrdered(left, right) >= 0
	case "<=":
		return compareOrdered(left, right) <= 0
	default:
		return false
	}
}

// compareOrdered compares ordered values (numbers, strings)
func compareOrdered(left, right interface{}) int {
	switch l := left.(type) {
	case int:
		if r, ok := right.(int); ok {
			if l < r {
				return -1
			} else if l > r {
				return 1
			}
			return 0
		}
	case string:
		if r, ok := right.(string); ok {
			if l < r {
				return -1
			} else if l > r {
				return 1
			}
			return 0
		}
	}
	return 0
}

// evaluateExpression evaluates an expression to a value
func (db *Database) evaluateExpression(expr parser.Expression) (interface{}, error) {
	switch e := expr.(type) {
	case *parser.Literal:
		return e.Value, nil
	case *parser.Identifier:
		return nil, fmt.Errorf("identifiers not supported in value context")
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// ResultSet represents the result of a SELECT query
type ResultSet struct {
	Columns []string
	Rows    [][]interface{}
}

// Print prints the result set in a formatted way
func (rs *ResultSet) Print() {
	if len(rs.Rows) == 0 {
		fmt.Println("No results")
		return
	}

	// Print column headers
	for i, col := range rs.Columns {
		if i > 0 {
			fmt.Print("\t")
		}
		fmt.Print(col)
	}
	fmt.Println()

	// Print separator
	for i := range rs.Columns {
		if i > 0 {
			fmt.Print("\t")
		}
		fmt.Print("--------")
	}
	fmt.Println()

	// Print rows
	for _, row := range rs.Rows {
		for i, value := range row {
			if i > 0 {
				fmt.Print("\t")
			}
			if value == nil {
				fmt.Print("NULL")
			} else {
				fmt.Print(value)
			}
		}
		fmt.Println()
	}
}
