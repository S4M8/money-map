.PHONY: help up down build rebuild logs ps clean nuke

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  up        - Start the application in detached mode"
	@echo "  down      - Stop the application"
	@echo "  build     - Build the application images"
	@echo "  rebuild   - Rebuild the application images from scratch"
	@echo "  logs      - View the application logs"
	@echo "  ps        - Show the status of the application containers"
	@echo "  clean     - Stop the application and remove all data"
	@echo "  nuke      - Stop the application, remove all data, and remove the app image"

up:
	@echo "Starting the application..."
	@docker-compose up -d

down:
	@echo "Stopping the application..."
	@docker-compose down

build:
	@echo "Building the application..."
	@docker-compose build

rebuild:
	@echo "Rebuilding the application from scratch..."
	@docker-compose build --no-cache

logs:
	@echo "Showing application logs..."
	@docker-compose logs -f

ps:
	@echo "Showing application status..."
	@docker-compose ps

clean:
	@echo "Stopping the application and removing all data..."
	@docker-compose down -v

nuke:
	@echo "Stopping the application, removing all data, and removing the app image..."
	@docker-compose down --volumes
	@docker rmi money-map_app || true
