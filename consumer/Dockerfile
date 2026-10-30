FROM golang:1.27.3-alpine AS build
WORKDIR /src
ARG GOPROXY=https://proxy.golang.org,direct
COPY go.mod go.sum ./
RUN --mount=type=secret,id=ca_bundle,target=/etc/ssl/certs/ca-certificates.crt go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /service ./cmd/outbox-consumer
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /service /service
USER 10001:10001
ENTRYPOINT ["/service"]
