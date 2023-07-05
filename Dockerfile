FROM golang:1.16-alpine AS builder 

WORKDIR /app

COPY main .
COPY . .
# COPY migrate.linux-amd64 . 

# FROM golang:1.16-alpine AS builder 

# # Set destination for COPY
# WORKDIR /app
# # Download Go modules
# COPY go.mod .
# RUN go mod download
# # Copy the source code.
# COPY . .g
# # Build
# # RUN apk add curl
# # RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz
# COPY migrate.linux-amd64 .
# RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go 

# RUN STAGE
FROM alpine:latest 
WORKDIR /app
COPY --from=builder /app/main .
# COPY --from=builder /app/migrate.linux-amd64  ./migrate
COPY app.env .
COPY src/db/migration ./src/db/migration 
COPY wait-for.sh .
COPY start.sh .


EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh"]