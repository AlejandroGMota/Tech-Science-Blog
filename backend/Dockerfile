# Pure-Go build: SQLite via modernc.org/sqlite, no cgo, no Oracle Instant
# Client. The Oracle store is only compiled with -tags oracle.
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/tech-blog-api ./cmd/api

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 1000 app
WORKDIR /app
COPY --from=build /out/tech-blog-api /app/tech-blog-api
# public-dist/ and admin-dist/ are built by the frontend and mounted at run
# time; the database lives in /app/data, a volume.
USER app
EXPOSE 8081
ENTRYPOINT ["/app/tech-blog-api"]
