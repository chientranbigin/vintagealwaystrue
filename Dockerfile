# Build stage
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

# Ensure go.mod is synchronized and build static binary
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Final stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

# Copy binary, static assets, templates and public files
COPY --from=builder /app/server /app/server
COPY --from=builder /app/static /app/static
COPY --from=builder /app/public /app/public
COPY --from=builder /app/templates /app/templates

ENV PORT=8080
EXPOSE 8080

CMD ["/app/server"]
