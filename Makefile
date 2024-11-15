up:
	docker-compose up -d

down:
	docker-compose down -v --rmi=local --remove-orphans

new-migration:
	cd database/migrations && sql-migrate new $(name)

migrations-up:
	cd database/migrations && sql-migrate up

migrations-down:
	cd database/migrations && sql-migrate down

build:
	docker build -t $(IMAGE_NAME) .

lint-dockerfile:
	hadolint Dockerfile
