package engine

import (
	"fmt"
	"go-rdbms/parser"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Storage handles file-based persistence of database tables
type Storage struct {
	dataDir string
}

// NewStorage creates a new storage instance
func NewStorage(dataDir string) *Storage {
	return &Storage{
		dataDir: dataDir,
	}
}

// Init initializes the storage directory
func (s *Storage) Init() error {
	return os.MkdirAll(s.dataDir, 0755)
}

// SaveTable saves a table to disk
func (s *Storage) SaveTable(table *Table) error {
	filename := s.getTableFilename(table.Name)
	csvData := table.ToCSV()

	return ioutil.WriteFile(filename, []byte(csvData), 0644)
}

// LoadTable loads a table from disk
func (s *Storage) LoadTable(tableName string) (*Table, error) {
	filename := s.getTableFilename(tableName)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("table %s does not exist on disk", tableName)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading table file: %v", err)
	}

	return TableFromCSV(tableName, string(data))
}

// DeleteTable removes a table file from disk
func (s *Storage) DeleteTable(tableName string) error {
	filename := s.getTableFilename(tableName)
	return os.Remove(filename)
}

// ListTables returns all table names stored on disk
func (s *Storage) ListTables() ([]string, error) {
	files, err := ioutil.ReadDir(s.dataDir)
	if err != nil {
		return nil, err
	}

	var tables []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".table") {
			tableName := strings.TrimSuffix(file.Name(), ".table")
			tables = append(tables, tableName)
		}
	}

	return tables, nil
}

// SaveDatabase saves all tables in a database
func (s *Storage) SaveDatabase(db *Database) error {
	for tableName, table := range db.Tables {
		if err := s.SaveTable(table); err != nil {
			return fmt.Errorf("error saving table %s: %v", tableName, err)
		}
	}
	return nil
}

// LoadDatabase loads all tables from disk into a database
func (s *Storage) LoadDatabase(db *Database) error {
	tableNames, err := s.ListTables()
	if err != nil {
		return err
	}

	for _, tableName := range tableNames {
		table, err := s.LoadTable(tableName)
		if err != nil {
			return fmt.Errorf("error loading table %s: %v", tableName, err)
		}
		db.Tables[tableName] = table
	}

	return nil
}

// getTableFilename returns the filename for a table
func (s *Storage) getTableFilename(tableName string) string {
	return filepath.Join(s.dataDir, tableName+".table")
}

// PersistedDatabase combines Database with Storage for automatic persistence
type PersistedDatabase struct {
	*Database
	storage *Storage
}

// NewPersistedDatabase creates a new database with automatic file persistence
func NewPersistedDatabase(dataDir string) (*PersistedDatabase, error) {
	storage := NewStorage(dataDir)
	if err := storage.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %v", err)
	}

	db := NewDatabase()
	pdb := &PersistedDatabase{
		Database: db,
		storage:  storage,
	}

	// Load existing tables
	if err := storage.LoadDatabase(db); err != nil {
		return nil, fmt.Errorf("failed to load database: %v", err)
	}

	return pdb, nil
}

// ExecuteCreateTable executes CREATE TABLE and saves to disk
func (pdb *PersistedDatabase) ExecuteCreateTable(stmt *parser.CreateTableStatement) error {
	if err := pdb.Database.ExecuteCreateTable(stmt); err != nil {
		return err
	}

	// Save the new table
	table := pdb.Tables[stmt.TableName]
	return pdb.storage.SaveTable(table)
}

// ExecuteInsert executes INSERT and saves to disk
func (pdb *PersistedDatabase) ExecuteInsert(stmt *parser.InsertStatement) error {
	if err := pdb.Database.ExecuteInsert(stmt); err != nil {
		return err
	}

	// Save the updated table
	table := pdb.Tables[stmt.TableName]
	return pdb.storage.SaveTable(table)
}

// ExecuteUpdate executes UPDATE and saves to disk
func (pdb *PersistedDatabase) ExecuteUpdate(stmt *parser.UpdateStatement) error {
	if err := pdb.Database.ExecuteUpdate(stmt); err != nil {
		return err
	}

	// Save the updated table
	table := pdb.Tables[stmt.TableName]
	return pdb.storage.SaveTable(table)
}

// ExecuteDelete executes DELETE and saves to disk
func (pdb *PersistedDatabase) ExecuteDelete(stmt *parser.DeleteStatement) error {
	if err := pdb.Database.ExecuteDelete(stmt); err != nil {
		return err
	}

	// Save the updated table
	table := pdb.Tables[stmt.TableName]
	return pdb.storage.SaveTable(table)
}

// ExecuteSelect executes SELECT (no persistence needed)
func (pdb *PersistedDatabase) ExecuteSelect(stmt *parser.SelectStatement) (*ResultSet, error) {
	return pdb.Database.ExecuteSelect(stmt)
}
