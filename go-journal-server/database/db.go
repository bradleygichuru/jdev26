package database

import (
	"fmt"
	"go-rdbms/engine"
	"go-rdbms/parser"
	"strconv"
	"strings"
	"time"
)

type JournalDB struct {
	db *engine.PersistedDatabase
}

type JournalEntryDB struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      string    `json:"tags"` // stored as comma-separated string
}

func NewJournalDB(dataDir string) (*JournalDB, error) {
	pdb, err := engine.NewPersistedDatabase(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	jdb := &JournalDB{db: pdb}

	// Initialize schema
	if err := jdb.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return jdb, nil
}

func (j *JournalDB) initSchema() error {
	// Create entries table if it doesn't exist
	createStmt := &parser.CreateTableStatement{
		TableName: "entries",
		Columns: []*parser.ColumnDefinition{
			{Name: "id", DataType: parser.DATATYPE_INTEGER, PrimaryKey: true},
			{Name: "title", DataType: parser.DATATYPE_TEXT},
			{Name: "content", DataType: parser.DATATYPE_TEXT},
			{Name: "created_at", DataType: parser.DATATYPE_TEXT},
			{Name: "updated_at", DataType: parser.DATATYPE_TEXT},
			{Name: "tags", DataType: parser.DATATYPE_TEXT},
		},
	}

	err := j.db.ExecuteCreateTable(createStmt)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}

	return nil
}

func (j *JournalDB) CreateEntry(title, content string, tags []string) (*JournalEntryDB, error) {
	now := time.Now()

	tagsStr := strings.Join(tags, ",")
	if tagsStr == "" {
		tagsStr = ","
	}

	// Get next ID
	nextID, err := j.getNextID()
	if err != nil {
		return nil, err
	}

	insertStmt := &parser.InsertStatement{
		TableName: "entries",
		Values: []parser.Expression{
			&parser.Literal{Value: nextID, Type: parser.DATATYPE_INTEGER},
			&parser.Literal{Value: title, Type: parser.DATATYPE_TEXT},
			&parser.Literal{Value: content, Type: parser.DATATYPE_TEXT},
			&parser.Literal{Value: now.Format(time.RFC3339), Type: parser.DATATYPE_TEXT},
			&parser.Literal{Value: now.Format(time.RFC3339), Type: parser.DATATYPE_TEXT},
			&parser.Literal{Value: tagsStr, Type: parser.DATATYPE_TEXT},
		},
	}

	err = j.db.ExecuteInsert(insertStmt)
	if err != nil {
		return nil, err
	}

	return &JournalEntryDB{
		ID:        nextID,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		Tags:      tagsStr,
	}, nil
}

func (j *JournalDB) GetEntry(id int) (*JournalEntryDB, error) {
	selectStmt := &parser.SelectStatement{
		TableName: "entries",
		Columns:   []parser.Expression{&parser.StarExpression{}},
		Where: &parser.BinaryExpression{
			Left:     &parser.Identifier{Value: "id"},
			Operator: "=",
			Right:    &parser.Literal{Value: id, Type: parser.DATATYPE_INTEGER},
		},
	}

	result, err := j.db.ExecuteSelect(selectStmt)
	if err != nil {
		return nil, err
	}

	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("entry not found")
	}

	return j.rowToEntry(result.Rows[0], result.Columns)
}

func (j *JournalDB) GetAllEntries() ([]*JournalEntryDB, error) {
	selectStmt := &parser.SelectStatement{
		TableName: "entries",
		Columns:   []parser.Expression{&parser.StarExpression{}},
	}

	result, err := j.db.ExecuteSelect(selectStmt)
	if err != nil {
		return nil, err
	}

	entries := make([]*JournalEntryDB, 0, len(result.Rows))
	for _, row := range result.Rows {
		entry, err := j.rowToEntry(row, result.Columns)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (j *JournalDB) SearchEntries(query string) ([]*JournalEntryDB, error) {
	selectStmt := &parser.SelectStatement{
		TableName: "entries",
		Columns:   []parser.Expression{&parser.StarExpression{}},
		Where: &parser.BinaryExpression{
			Left: &parser.BinaryExpression{
				Left:     &parser.Identifier{Value: "title"},
				Operator: "LIKE",
				Right:    &parser.Literal{Value: "%" + query + "%", Type: parser.DATATYPE_TEXT},
			},
			Operator: "OR",
			Right: &parser.BinaryExpression{
				Left:     &parser.Identifier{Value: "content"},
				Operator: "LIKE",
				Right:    &parser.Literal{Value: "%" + query + "%", Type: parser.DATATYPE_TEXT},
			},
		},
	}

	result, err := j.db.ExecuteSelect(selectStmt)
	if err != nil {
		return nil, err
	}

	entries := make([]*JournalEntryDB, 0, len(result.Rows))
	for _, row := range result.Rows {
		entry, err := j.rowToEntry(row, result.Columns)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (j *JournalDB) UpdateEntry(id int, title, content *string, tags []string) error {
	updates := make(map[string]parser.Expression)

	if title != nil {
		updates["title"] = &parser.Literal{Value: *title, Type: parser.DATATYPE_TEXT}
	}

	if content != nil {
		updates["content"] = &parser.Literal{Value: *content, Type: parser.DATATYPE_TEXT}
	}

	if tags != nil {
		tagsStr := strings.Join(tags, ",")
		if tagsStr == "" {
			tagsStr = ","
		}
		updates["tags"] = &parser.Literal{Value: tagsStr, Type: parser.DATATYPE_TEXT}
	}

	if len(updates) > 0 {
		updates["updated_at"] = &parser.Literal{Value: time.Now().Format(time.RFC3339), Type: parser.DATATYPE_TEXT}

		updateStmt := &parser.UpdateStatement{
			TableName: "entries",
			Set:       updates,
			Where: &parser.BinaryExpression{
				Left:     &parser.Identifier{Value: "id"},
				Operator: "=",
				Right:    &parser.Literal{Value: id, Type: parser.DATATYPE_INTEGER},
			},
		}

		return j.db.ExecuteUpdate(updateStmt)
	}

	return nil
}

func (j *JournalDB) DeleteEntry(id int) error {
	deleteStmt := &parser.DeleteStatement{
		TableName: "entries",
		Where: &parser.BinaryExpression{
			Left:     &parser.Identifier{Value: "id"},
			Operator: "=",
			Right:    &parser.Literal{Value: id, Type: parser.DATATYPE_INTEGER},
		},
	}

	return j.db.ExecuteDelete(deleteStmt)
}

func (j *JournalDB) getNextID() (int, error) {
	selectStmt := &parser.SelectStatement{
		TableName: "entries",
		Columns:   []parser.Expression{&parser.Identifier{Value: "id"}},
	}

	result, err := j.db.ExecuteSelect(selectStmt)
	if err != nil {
		return 1, nil // If table is empty, start with 1
	}

	maxID := 0
	for _, row := range result.Rows {
		// Create a map from column names to values
		data := make(map[string]interface{})
		for i, col := range result.Columns {
			if i < len(row) {
				data[col] = row[i]
			}
		}

		if idVal, ok := data["id"]; ok {
			switch v := idVal.(type) {
			case int:
				if v > maxID {
					maxID = v
				}
			case string:
				if id, err := strconv.Atoi(v); err == nil && id > maxID {
					maxID = id
				}
			}
		}
	}

	return maxID + 1, nil
}

func (j *JournalDB) rowToEntry(row []interface{}, columns []string) (*JournalEntryDB, error) {
	entry := &JournalEntryDB{}

	// Create a map from column names to values
	data := make(map[string]interface{})
	for i, col := range columns {
		if i < len(row) {
			data[col] = row[i]
		}
	}

	// Handle ID
	if idVal, ok := data["id"]; ok {
		switch v := idVal.(type) {
		case int:
			entry.ID = v
		case string:
			if id, err := strconv.Atoi(v); err == nil {
				entry.ID = id
			}
		}
	}

	// Handle strings
	if title, ok := data["title"].(string); ok {
		entry.Title = title
	}
	if content, ok := data["content"].(string); ok {
		entry.Content = content
	}
	if tags, ok := data["tags"].(string); ok {
		tags = strings.TrimSpace(tags)
		if tags == "," {
			tags = ""
		}
		entry.Tags = tags
	}

	// Handle timestamps
	if createdAt, ok := data["created_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			entry.CreatedAt = t
		}
	}
	if updatedAt, ok := data["updated_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
			entry.UpdatedAt = t
		}
	}

	return entry, nil
}
