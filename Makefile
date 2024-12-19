help:
	@echo "Usage:"
	@echo ""
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

up: ## Start the application
	docker compose up -d

down: ## Stop the application
	docker compose down -v --rmi=local --remove-orphans

new-migration: ## Create a new migration
	cd database/migrations && sql-migrate new $(name)

migrations-up: ## Run all migrations
	cd database/migrations && sql-migrate up

migrations-down: ## Rollback all migrations
	cd database/migrations && sql-migrate down

build: ## Build the application
	docker build -t $(IMAGE_NAME) .

lint-dockerfile: ## Lint Dockerfile
	hadolint Dockerfile
