# Stage 1: Build
FROM golang:1.20-alpine AS build
WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp

# Stage 2: Run
FROM alpine:latest
WORKDIR /app

COPY --from=build /app/myapp /app/myapp

EXPOSE 8080
CMD ["/app/myapp"]
