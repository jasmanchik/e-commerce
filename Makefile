up: docker-build docker-up

down:
	docker-compose down

docker-build:
	docker-compose build

docker-up:
	args=$(args) docker-compose up

