install:
	go mod tidy
	go install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

protoc:
	protoc -I .  \
	-I$$(go list -m -f "{{.Dir}}" github.com/grpc-ecosystem/grpc-gateway/v2) \
	-I$$(go list -m -f "{{.Dir}}" github.com/grpc-ecosystem/grpc-gateway/v2)/third_party/googleapis \
	-I$$(go list -m -f "{{.Dir}}" github.com/grpc-ecosystem/grpc-gateway/v2)/protoc-gen-openapiv2 \
	--go_out . \
	--go_opt paths=source_relative \
	--go-grpc_out . \
	--go-grpc_opt paths=source_relative \
	--grpc-gateway_out . \
    --grpc-gateway_opt paths=source_relative \
	--openapiv2_out . \
	./service/pb/*.proto
	mv service/pb/wikitable.swagger.json swagger/

test:
	go test -v -cover -race -p 1 ./...

build:
	docker build -t wikitable-api .

run: build
	docker run --rm -p 8080:8080 -e PORT=8080 wikitable-api

cover-profile:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm -f coverage.out