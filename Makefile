all: run

.SILENT: run
run:
	go run ./cmd/run

.SILENT: test
test:
	go test ./...
