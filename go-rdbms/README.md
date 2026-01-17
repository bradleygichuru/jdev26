# Simple RDBMS

A simple relational database management system implemented in Go, supporting basic CRUD operations, primary keys, unique constraints, and JOINs.

## Features

- **Data Types**: INTEGER, TEXT, BOOLEAN
- **CRUD Operations**: CREATE TABLE, INSERT, SELECT, UPDATE, DELETE
- **Constraints**: PRIMARY KEY, UNIQUE
- **Queries**: Basic SELECT with WHERE conditions, INNER JOIN
- **Storage**: File-based persistence with CSV-like format
- **Interface**: Interactive REPL with SQL commands

## Usage

### Building

```bash
go build -o rdbms .
```

### Running

```bash
./rdbms
```

This starts an interactive REPL where you can enter SQL commands.

### Example Session

```
rdbms> CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER);
Table users created successfully

rdbms> INSERT INTO users VALUES (1, 'Alice', 25);
Row inserted successfully

rdbms> INSERT INTO users VALUES (2, 'Bob', 30);
Row inserted successfully

rdbms> SELECT * FROM users;
id	name	age
--------	--------	--------
1	Alice	25
2	Bob	30

rdbms> SELECT * FROM users WHERE age > 25;
id	name	age
--------	--------	--------
2	Bob	30

rdbms> CREATE TABLE posts (id INTEGER PRIMARY KEY, user_id INTEGER, title TEXT);
Table posts created successfully

rdbms> INSERT INTO posts VALUES (1, 1, 'Hello World');
Row inserted successfully

rdbms> SELECT users.name, posts.title FROM users JOIN posts ON users.id = posts.user_id;
name	title
--------	--------
Alice	Hello World

rdbms> exit
Goodbye!
```

## SQL Syntax

### CREATE TABLE
```sql
CREATE TABLE table_name (
    column1 datatype [PRIMARY KEY],
    column2 datatype [UNIQUE],
    ...
);
```

### INSERT
```sql
INSERT INTO table_name VALUES (value1, value2, ...);
```

### SELECT
```sql
SELECT * FROM table_name [WHERE condition] [JOIN other_table ON condition];
SELECT column1, column2 FROM table_name [WHERE condition];
```

### UPDATE
```sql
UPDATE table_name SET column1 = value1, column2 = value2 WHERE condition;
```

### DELETE
```sql
DELETE FROM table_name WHERE condition;
```

## Architecture

- **Parser**: Recursive descent SQL parser with lexer
- **Engine**: In-memory database with file persistence
- **Storage**: CSV-based file storage with schema headers
- **REPL**: Interactive command-line interface

## Limitations

- Single-table WHERE conditions (no complex expressions)
- Equality JOINs only
- No transactions
- No indexes beyond primary key
- No aggregate functions (SUM, COUNT, etc.)
- Limited error recovery

## Files

- `main.go`: Entry point and REPL initialization
- `repl/repl.go`: Interactive REPL implementation
- `parser/lexer.go`: SQL lexical analysis
- `parser/parser.go`: SQL parsing
- `parser/ast.go`: Abstract syntax tree definitions
- `engine/database.go`: Database operations
- `engine/table.go`: Table and row management
- `engine/storage.go`: File-based persistence
- `examples/sample.sql`: Sample SQL commands