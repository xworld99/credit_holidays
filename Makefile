up:
	docker compose -f ./deployments/docker-compose.yml up -d --build

down:
	docker compose -f ./deployments/docker-compose.yml down

make swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --parseDependency --parseInternal -g cmd/main.go

test:
	echo "Not implemented"

wire:
	go install github.com/google/wire/cmd/wire@latest
	wire credit_holidays/internal/handlers/.


