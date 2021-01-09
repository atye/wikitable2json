FROM golang:1.14 as staging
RUN apt update -y && apt install -y curl unzip ca-certificates && update-ca-certificates

RUN curl -L https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip -O
RUN unzip -o protoc-3.14.0-linux-x86_64.zip -d /usr/local bin/protoc && unzip -o protoc-3.14.0-linux-x86_64.zip -d /usr/local 'include/*'
RUn chmod +x /usr/local/bin/protoc

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
