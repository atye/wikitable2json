FROM golang:1.18.3 as staging
RUN apt update -y && apt install -y curl unzip ca-certificates && update-ca-certificates

WORKDIR /wikitable2json
COPY . .
RUN CGO_ENABLED=0 go build -o /tmp/service ./cmd/main.go && chmod u+x /tmp/service

FROM scratch
COPY --from=staging /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=staging /tmp/service /service
ENTRYPOINT ["/service"]
