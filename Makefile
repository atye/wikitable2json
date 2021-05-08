test:
	go test -count=1 -v -cover -race ./...

build:
	docker build -t wikitable-api .

run: build
	docker run --rm -p 8080:8080 -e PORT=8080 wikitable-api

cover-profile:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm -f coverage.out
