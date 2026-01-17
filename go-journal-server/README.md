# Go Journal Server

A simple REST API for managing journal entries built with Go, Chi router, and a custom RDBMS.

## Features

- Create, read, update, and delete journal entries
- Search entries by title, content, or tags
- Persistent storage using custom RDBMS
- RESTful API with JSON responses

## API Endpoints

### Health Check
- `GET /` - Server health check

### Journal Entries

#### Create Entry
- `POST /api/entries`
- Body: `{"title": "string", "content": "string", "tags": ["string"]}`

#### Get All Entries
- `GET /api/entries`

#### Get Entry by ID
- `GET /api/entries/{id}`

#### Update Entry
- `PUT /api/entries/{id}`
- Body: `{"title": "string", "content": "string", "tags": ["string"]}` (partial updates supported)

#### Delete Entry
- `DELETE /api/entries/{id}`

#### Search Entries
- `GET /api/entries/search?q={query}`
- Searches across title, content, and tags

## Response Format

All responses follow this format:
```json
{
  "success": true|false,
  "data": {...},
  "error": "error message"
}
```

## Database

Uses a custom RDBMS with the following schema:
```sql
CREATE TABLE entries (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    tags TEXT
);
```

## Running the Server

```bash
go build -o journal-server
./journal-server
```

Server runs on port 8080 by default.

## Dependencies

- `github.com/go-chi/chi/v5` - HTTP router
- Custom `go-rdbms` - Database implementation