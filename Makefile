up:
	docker-compose up -d

down:
	docker-compose down -v --rmi=local --remove-orphans
