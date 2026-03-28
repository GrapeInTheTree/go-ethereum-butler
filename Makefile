.PHONY: build test vet clean run

build:
	go build -o butler ./cmd/butler

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f butler

run:
	go run ./cmd/butler
