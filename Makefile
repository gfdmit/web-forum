up:
	docker compose up -d

build:
	docker compose up -d --build

restart:
	docker compose down
	docker compose up -d

rebuild:
	docker compose down
	docker compose up -d --build

rebuild_clear:
	docker compose down -v
	docker compose up -d --build

down:
	docker compose down

down_clear:
	docker compose down -v
