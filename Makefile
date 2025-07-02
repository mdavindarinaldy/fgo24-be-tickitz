include .env

DB_SSLMODE?=disable

MIGRATION_DIR?=./migrations

MIGRATE=migrate -source file://$(MIGRATION_DIR)	\
	-database "postgres://$(PGUSER):$(PGPASSWORD)@$(PGHOST):$(PGPORT)/$(PGNAME)?sslmode=$(DB_SSLMODE)"

migration_create:
	migrate create -seq -dir $(MIGRATION_DIR) -ext sql $(name)

migration_up:
	$(MIGRATE) up 1

migration_down:
	$(MIGRATE) down 1

migration_drop:
	$(MIGRATE) drop p -f

migration_force:
	$(MIGRATE) force $(ver)

build_win:
	GOOS=windows GOARC=amd64 go build -o build/app.exe main.go
build_linux:
	GOOS=linux GOARC=amd64 go build -o build/app main.go
build_darwin:
	GOOS=darwin GOARC=arm64 go build -o build/app_darwin main.go
clean:
	rm -rf build
build:
	clean build_win build_linux build_darwin