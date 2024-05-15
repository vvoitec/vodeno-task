sqlc-gen:
	sqlc generate

up: down
	docker-compose up

down:
	docker-compose down --remove-orphans
