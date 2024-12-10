test:
	go test -count=1 -race ./...

build:
	docker build -t wikitable2json .

run: build
	docker run --rm -p 8080:8080 -e PORT=8080 wikitable2json
