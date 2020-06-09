FROM golang:1.14 as staging
RUN apt update -y && apt install -y curl unzip ca-certificates && update-ca-certificates

WORKDIR /usr/local
RUN curl -Lk https://github.com/protocolbuffers/protobuf/releases/download/v3.12.3/protoc-3.12.3-linux-x86_64.zip -o protobuf.zip && \
    unzip protobuf.zip && \
    rm protobuf.zip

WORKDIR /wikitable-api
COPY . .
RUN make install protoc

RUN CGO_ENABLED=0 go build -o /tmp/service ./cmd/main.go && chmod u+x /tmp/service

FROM scratch
COPY --from=staging /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=staging /wikitable-api/static /static
COPY --from=staging /wikitable-api/swagger /swagger
COPY --from=staging /tmp/service /service
ENTRYPOINT ["/service"]
