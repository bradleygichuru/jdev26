module go-journal-server

go 1.25.5

require (
	github.com/go-chi/chi/v5 v5.2.4
	go-rdbms v0.0.0
)

require github.com/go-chi/cors v1.2.2 // indirect

replace go-rdbms => ../go-rdbms
