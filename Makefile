up:
	docker compose --env-file .env up -d

build:
	docker compose build

up_build:
	docker compose --env-file .env up -d --build

stop:
	docker compose stop

down:
	docker compose down
