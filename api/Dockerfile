FROM golang:1.27.3-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=secret,id=ca_bundle,target=/etc/ssl/certs/ca-certificates.crt go mod download
COPY .. .
RUN CGO_ENABLED=0 go build -trimpath -o /api ./cmd/api

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
USER 10001:10001
COPY --from=build /api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
