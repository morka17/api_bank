# syntax=docker/dockerfile:1
FROM golang:11.16-alpine AS builder 

# Set destination for COPY
WORKDIR /app
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code.
COPY . .
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go 


# RUN STAGE
FROM alpine:latest 
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080

# Run
CMD ["/docker-gs-ping"]