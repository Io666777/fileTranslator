up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

reup:
	make down
	make build
	make up

