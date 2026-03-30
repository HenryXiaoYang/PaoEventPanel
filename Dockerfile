# -- Build frontend --
FROM node:22-alpine AS frontend
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# -- Build backend --
FROM golang:1.25-alpine AS backend
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN CGO_ENABLED=1 go build -o pao-event-panel .

# -- Runtime --
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend /app/pao-event-panel .
COPY --from=backend /app/web/dist ./web/dist
RUN mkdir -p uploads
EXPOSE 8080
CMD ["./pao-event-panel"]
