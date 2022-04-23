test:
	go test -coverpkg=./... -count=1 -race ./...

build:
	docker build -t wikitable2json .

run: build
	docker run --rm -p 8080:8080 -e PORT=8080 wikitable2json

cover:
	cd internal/entrypoint && go test -coverpkg=../... -coverprofile=coverage.out ./...
	go tool cover -html=internal/entrypoint/coverage.out
	rm -f internal/entrypoint/coverage.out
