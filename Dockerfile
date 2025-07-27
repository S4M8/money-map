# Stage 1: Build the React frontend
FROM node:20-alpine AS builder
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm install
COPY web/. .
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.24-alpine
WORKDIR /app
COPY --from=builder /app/web/build ./web/build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
RUN go build -o /money-map ./cmd/server

EXPOSE 8080
CMD [ "/money-map" ]