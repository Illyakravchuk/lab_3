tests:
	go test ./...

out/example: cmd/painter/main.go
	mkdir -p out
	go build -o out/example ./cmd/painter