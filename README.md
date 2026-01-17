# Journal Application

A full-stack journal application with a custom relational database management system. This project consists of three main components: a Go backend API, a Svelte frontend, and a custom RDBMS engine.

## Architecture

- **go-journal-server/** - REST API backend built with Go and Chi router
- **svelte-frontend/** - Modern web interface built with SvelteKit and Tailwind CSS
- **go-rdbms/** - Custom relational database management system with SQL parser

The journal server uses the custom RDBMS for data persistence, providing a complete full-stack solution.

## Prerequisites

- **Go** 1.25 or later
- **Node.js** 18 or later


## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd jdev26
```

### 2. Start the Backend Server

Navigate to the Go server directory and build:

```bash
cd go-journal-server
go build -o journal-server
./journal-server
```

The server will start on port 8080.

### 3. Start the Frontend Development Server

In a new terminal, navigate to the frontend directory:

```bash
cd svelte-frontend
npm install
npm run dev
```

The frontend will typically start on http://localhost:5173

### 4. Run Both Services for Testing

To test the full application, keep both services running:
- Backend: `./go-journal-server/journal-server` (port 8080)
- Frontend: `cd svelte-frontend && npm run dev` (port 5173)

## Project Structure

```
jdev26/
├── go-journal-server/     # Backend REST API
│   ├── handlers/          # HTTP handlers and routes
│   ├── database/          # Database integration
│   └── main.go           # Server entry point
├── svelte-frontend/       # Frontend web application
│   ├── src/
│   │   ├── routes/       # SvelteKit routes
│   │   ├── lib/          # Components and utilities
│   │   └── app.html      # App template
│   └── package.json      # Frontend dependencies
└── go-rdbms/             # Custom RDBMS engine
    ├── parser/           # SQL parser components
    ├── engine/           # Database engine
    └── repl/             # Interactive command-line interface
```

## Development Notes

### AI Assistance Attribution

The following components in the `go-rdbms/` module were developed with AI assistance:

- `parser/lexer.go`: SQL lexical analysis
- `parser/parser.go`: SQL parsing  
- `parser/ast.go`: Abstract syntax tree definitions

These components provide the core SQL parsing functionality for the custom database engine.

## Component Details

### Backend (go-journal-server)

Provides REST API endpoints for managing journal entries:
- CRUD operations for entries
- Search functionality
- JSON API responses
- Custom RDBMS integration

See [go-journal-server/README.md](go-journal-server/README.md) for detailed API documentation.

### Frontend (svelte-frontend)

Modern web interface with:
- SvelteKit framework
- Tailwind CSS styling
- TypeScript support
- Component-based UI library

See [svelte-frontend/README.md](svelte-frontend/README.md) for frontend-specific setup.

### Database Engine (go-rdbms)

Custom relational database with:
- SQL parser and lexer
- Support for basic SQL operations
- File-based persistence
- Interactive REPL interface

See [go-rdbms/README.md](go-rdbms/README.md) for complete database documentation.

## Troubleshooting

### Backend Issues
- Ensure Go 1.25+ is installed
- Check that the custom RDBMS dependency path is correct
- Verify port 8080 is available

### Frontend Issues
- Ensure Node.js 18+ is installed
- Run `npm install` to update dependencies
- Check that the backend is running on port 8080

### Database Issues
- The custom RDBMS stores data in local files
- Check file permissions in the data directories
- Use the REPL interface for direct database testing

