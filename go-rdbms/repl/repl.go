package repl

import (
	"bufio"
	"fmt"
	"go-rdbms/engine"
	"go-rdbms/parser"
	"os"
	"strings"
)

// Repl represents the interactive read-eval-print loop
type Repl struct {
	database *engine.PersistedDatabase
}

// NewRepl creates a new REPL instance
func NewRepl(dataDir string) (*Repl, error) {
	db, err := engine.NewPersistedDatabase(dataDir)
	if err != nil {
		return nil, err
	}

	return &Repl{
		database: db,
	}, nil
}

// Start begins the interactive REPL session
func (r *Repl) Start() error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("rdbms> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if err := r.handleCommand(input); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	return scanner.Err()
}

// handleCommand processes a single command
func (r *Repl) handleCommand(input string) error {
	switch strings.ToLower(input) {
	case "exit", "quit", "\\q":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "help", "\\h", "?":
		r.showHelp()
	case "tables":
		r.showTables()
	default:
		return r.executeSQL(input)
	}
	return nil
}

// executeSQL parses and executes SQL commands
func (r *Repl) executeSQL(sql string) error {
	// Split SQL by semicolons and execute each statement
	statements := strings.Split(sql, ";")
	for _, stmtSQL := range statements {
		stmtSQL = strings.TrimSpace(stmtSQL)
		if stmtSQL == "" {
			continue
		}

		lexer := parser.NewLexer(stmtSQL)
		p := parser.NewParser(lexer)

		stmt, err := p.ParseStatement()
		if err != nil {
			return fmt.Errorf("parse error: %v", err)
		}

		if len(p.GetErrors()) > 0 {
			return fmt.Errorf("parse errors: %v", strings.Join(p.GetErrors(), "; "))
		}

		switch s := stmt.(type) {
		case *parser.CreateTableStatement:
			err = r.database.ExecuteCreateTable(s)
			if err == nil {
				fmt.Printf("Table %s created successfully\n", s.TableName)
			}
		case *parser.InsertStatement:
			err = r.database.ExecuteInsert(s)
			if err == nil {
				fmt.Println("Row inserted successfully")
			}
		case *parser.SelectStatement:
			result, execErr := r.database.ExecuteSelect(s)
			err = execErr
			if err == nil {
				result.Print()
			}
		case *parser.UpdateStatement:
			err = r.database.ExecuteUpdate(s)
			if err == nil {
				fmt.Println("Rows updated successfully")
			}
		case *parser.DeleteStatement:
			err = r.database.ExecuteDelete(s)
			if err == nil {
				fmt.Println("Rows deleted successfully")
			}
		default:
			return fmt.Errorf("unsupported statement type: %T", stmt)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// showTables displays all tables in the database
func (r *Repl) showTables() {
	if len(r.database.Tables) == 0 {
		fmt.Println("No tables found")
		return
	}

	fmt.Println("Tables:")
	for tableName := range r.database.Tables {
		fmt.Printf("  %s\n", tableName)
	}
}

// showHelp displays available commands
func (r *Repl) showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help, \\h, ?     - Show this help")
	fmt.Println("  exit, quit, \\q  - Exit the REPL")
	fmt.Println("  SQL commands coming soon...")
}
