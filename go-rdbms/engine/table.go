package engine

import (
	"fmt"
	"go-rdbms/parser"
	"strconv"
	"strings"
)

// Database represents the main database instance
type Database struct {
	Tables map[string]*Table
}

// NewDatabase creates a new database instance
func NewDatabase() *Database {
	return &Database{
		Tables: make(map[string]*Table),
	}
}

// Table represents a database table
type Table struct {
	Name       string
	Columns    []*Column
	Rows       []*Row
	PrimaryKey string               // column name of primary key
	index      map[interface{}]*Row // simple hash index for primary key
}

// NewTable creates a new table with the given schema
func NewTable(name string, columns []*Column) *Table {
	table := &Table{
		Name:    name,
		Columns: columns,
		Rows:    []*Row{},
		index:   make(map[interface{}]*Row),
	}

	// Find primary key column
	for _, col := range columns {
		if col.PrimaryKey {
			table.PrimaryKey = col.Name
			break
		}
	}

	return table
}

// Column represents a table column
type Column struct {
	Name       string
	DataType   parser.DataType
	PrimaryKey bool
	Unique     bool
}

// Row represents a table row
type Row struct {
	Data map[string]interface{}
}

// NewRow creates a new row with empty data
func NewRow() *Row {
	return &Row{
		Data: make(map[string]interface{}),
	}
}

// GetValue gets a value from the row by column name
func (r *Row) GetValue(columnName string) interface{} {
	return r.Data[columnName]
}

// SetValue sets a value in the row by column name
func (r *Row) SetValue(columnName string, value interface{}) {
	r.Data[columnName] = value
}

// ValidateRow validates that a row conforms to table schema
func (t *Table) ValidateRow(row *Row) error {
	// Check all required columns are present
	for _, col := range t.Columns {
		if _, exists := row.Data[col.Name]; !exists {
			return fmt.Errorf("missing value for column %s", col.Name)
		}

		// Type validation
		if err := t.validateValueType(col, row.Data[col.Name]); err != nil {
			return err
		}
	}

	// Check primary key uniqueness
	if t.PrimaryKey != "" {
		if pkValue, exists := row.Data[t.PrimaryKey]; exists {
			if _, exists := t.index[pkValue]; exists {
				return fmt.Errorf("primary key violation: %v already exists", pkValue)
			}
		}
	}

	// Check unique constraints
	for _, col := range t.Columns {
		if col.Unique && !col.PrimaryKey {
			if value, exists := row.Data[col.Name]; exists {
				for _, existingRow := range t.Rows {
					if existingValue := existingRow.GetValue(col.Name); existingValue != nil && existingValue == value {
						return fmt.Errorf("unique constraint violation for column %s: %v already exists", col.Name, value)
					}
				}
			}
		}
	}

	return nil
}

// validateValueType validates that a value matches the expected type
func (t *Table) validateValueType(col *Column, value interface{}) error {
	switch col.DataType {
	case parser.DATATYPE_INTEGER:
		if _, ok := value.(int); !ok {
			return fmt.Errorf("column %s expects INTEGER, got %T", col.Name, value)
		}
	case parser.DATATYPE_TEXT:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("column %s expects TEXT, got %T", col.Name, value)
		}
	case parser.DATATYPE_BOOLEAN:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("column %s expects BOOLEAN, got %T", col.Name, value)
		}
	}
	return nil
}

// InsertRow inserts a row into the table
func (t *Table) InsertRow(row *Row) error {
	if err := t.ValidateRow(row); err != nil {
		return err
	}

	t.Rows = append(t.Rows, row)

	// Update index if primary key exists
	if t.PrimaryKey != "" {
		if pkValue, exists := row.Data[t.PrimaryKey]; exists {
			t.index[pkValue] = row
		}
	}

	return nil
}

// FindRowByPrimaryKey finds a row by primary key value
func (t *Table) FindRowByPrimaryKey(pkValue interface{}) *Row {
	if t.PrimaryKey == "" {
		return nil
	}
	return t.index[pkValue]
}

// UpdateRow updates a row by primary key
func (t *Table) UpdateRow(pkValue interface{}, updates map[string]interface{}) error {
	row := t.FindRowByPrimaryKey(pkValue)
	if row == nil {
		return fmt.Errorf("row with primary key %v not found", pkValue)
	}

	// Apply updates
	for colName, value := range updates {
		col := t.findColumn(colName)
		if col == nil {
			return fmt.Errorf("column %s does not exist", colName)
		}

		if err := t.validateValueType(col, value); err != nil {
			return err
		}

		row.SetValue(colName, value)
	}

	// Re-validate unique constraints
	for _, col := range t.Columns {
		if col.Unique && !col.PrimaryKey {
			value := row.GetValue(col.Name)
			for _, existingRow := range t.Rows {
				if existingRow != row {
					if existingValue := existingRow.GetValue(col.Name); existingValue != nil && existingValue == value {
						return fmt.Errorf("unique constraint violation for column %s: %v already exists", col.Name, value)
					}
				}
			}
		}
	}

	return nil
}

// DeleteRow deletes a row by primary key
func (t *Table) DeleteRow(pkValue interface{}) error {
	if t.PrimaryKey == "" {
		return fmt.Errorf("table has no primary key")
	}

	row := t.FindRowByPrimaryKey(pkValue)
	if row == nil {
		return fmt.Errorf("row with primary key %v not found", pkValue)
	}

	// Remove from rows slice
	for i, r := range t.Rows {
		if r == row {
			t.Rows = append(t.Rows[:i], t.Rows[i+1:]...)
			break
		}
	}

	// Remove from index
	delete(t.index, pkValue)

	return nil
}

// findColumn finds a column by name
func (t *Table) findColumn(name string) *Column {
	for _, col := range t.Columns {
		if col.Name == name {
			return col
		}
	}
	return nil
}

// SelectRows performs a basic select operation
func (t *Table) SelectRows(columns []string, whereCondition func(*Row) bool) []*Row {
	var result []*Row

	for _, row := range t.Rows {
		if whereCondition == nil || whereCondition(row) {
			result = append(result, row)
		}
	}

	return result
}

// GetColumnNames returns all column names in order
func (t *Table) GetColumnNames() []string {
	var names []string
	for _, col := range t.Columns {
		names = append(names, col.Name)
	}
	return names
}

// ToCSV converts table to CSV format (for storage)
func (t *Table) ToCSV() string {
	var lines []string

	// Schema header
	var schemaParts []string
	for _, col := range t.Columns {
		colDef := col.Name + ":" + col.DataType.String()
		if col.PrimaryKey {
			colDef += ":PRIMARY_KEY"
		}
		if col.Unique {
			colDef += ":UNIQUE"
		}
		schemaParts = append(schemaParts, colDef)
	}
	lines = append(lines, "# SCHEMA: "+strings.Join(schemaParts, ","))

	// Data rows
	colNames := t.GetColumnNames()
	for _, row := range t.Rows {
		var values []string
		for _, colName := range colNames {
			value := row.GetValue(colName)
			if value == nil {
				values = append(values, "")
			} else {
				values = append(values, formatValue(value))
			}
		}
		lines = append(lines, strings.Join(values, ","))
	}

	return strings.Join(lines, "\n")
}

// formatValue formats a value for CSV storage
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		// Escape quotes and wrap in quotes if contains comma or quote
		if strings.Contains(v, ",") || strings.Contains(v, "\"") {
			v = strings.ReplaceAll(v, "\"", "\"\"")
			return "\"" + v + "\""
		}
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// FromCSV loads table from CSV format
func TableFromCSV(name, csvData string) (*Table, error) {
	lines := strings.Split(csvData, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty CSV data")
	}

	// Parse schema
	schemaLine := lines[0]
	if !strings.HasPrefix(schemaLine, "# SCHEMA: ") {
		return nil, fmt.Errorf("invalid schema line: %s", schemaLine)
	}

	schemaStr := strings.TrimPrefix(schemaLine, "# SCHEMA: ")
	colDefs := strings.Split(schemaStr, ",")

	var columns []*Column
	for _, colDef := range colDefs {
		colDef = strings.TrimSpace(colDef)
		parts := strings.Split(colDef, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid column definition: %s", colDef)
		}

		col := &Column{
			Name:     parts[0],
			DataType: parseDataType(parts[1]),
		}

		for i := 2; i < len(parts); i++ {
			switch parts[i] {
			case "PRIMARY_KEY":
				col.PrimaryKey = true
			case "UNIQUE":
				col.Unique = true
			}
		}

		columns = append(columns, col)
	}

	table := NewTable(name, columns)

	// Parse data rows
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		values := parseCSVLine(line)
		if len(values) != len(columns) {
			return nil, fmt.Errorf("row %d has %d values, expected %d", i, len(values), len(columns))
		}

		row := NewRow()
		for j, value := range values {
			col := columns[j]
			parsedValue, err := parseValue(value, col.DataType)
			if err != nil {
				return nil, fmt.Errorf("error parsing value for column %s: %v", col.Name, err)
			}
			row.SetValue(col.Name, parsedValue)
		}

		if err := table.InsertRow(row); err != nil {
			return nil, fmt.Errorf("error inserting row %d: %v", i, err)
		}
	}

	return table, nil
}

// Helper functions
func parseDataType(s string) parser.DataType {
	switch s {
	case "INTEGER":
		return parser.DATATYPE_INTEGER
	case "TEXT":
		return parser.DATATYPE_TEXT
	case "BOOLEAN":
		return parser.DATATYPE_BOOLEAN
	default:
		return parser.DATATYPE_TEXT
	}
}

func parseValue(s string, dataType parser.DataType) (interface{}, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	switch dataType {
	case parser.DATATYPE_INTEGER:
		return strconv.Atoi(s)
	case parser.DATATYPE_BOOLEAN:
		switch strings.ToLower(s) {
		case "true", "1":
			return true, nil
		case "false", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("invalid boolean value: %s", s)
		}
	case parser.DATATYPE_TEXT:
		// Remove quotes if present
		if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
			s = s[1 : len(s)-1]
			s = strings.ReplaceAll(s, "\"\"", "\"")
		}
		return s, nil
	default:
		return s, nil
	}
}

func parseCSVLine(line string) []string {
	var values []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(line); i++ {
		char := line[i]

		switch {
		case char == '"' && !inQuotes:
			inQuotes = true
		case char == '"' && inQuotes && i+1 < len(line) && line[i+1] == '"':
			// Escaped quote
			current.WriteByte('"')
			i++ // Skip next quote
		case char == '"' && inQuotes:
			inQuotes = false
		case char == ',' && !inQuotes:
			values = append(values, current.String())
			current.Reset()
		default:
			current.WriteByte(char)
		}
	}

	values = append(values, current.String())
	return values
}
