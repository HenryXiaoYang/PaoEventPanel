.PHONY: dev-backend dev-frontend build-frontend build run

dev-backend:
	go run main.go

dev-frontend:
	cd web && npm run dev

build-frontend:
	cd web && npm run build

build: build-frontend
	go build -o pao-event-panel .

run: build
	./pao-event-panel
