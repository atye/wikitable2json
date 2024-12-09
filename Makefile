test:
	go test -count=1 -race ./...

build:
	docker build -t wikitable2json .

run: build
	docker run --rm -p 8080:8080 -e PORT=8080 -e CACHE_SIZE=10 -e CACHE_EXPIRATION=10s wikitable2json
