up:
	docker compose -f ./deployments/docker-compose.yml up -d --build

down:
	docker compose -f ./deployments/docker-compose.yml down

test:
	echo "Not implemented"

wire:
	go install github.com/google/wire/cmd/wire@latest
	wire credit_holidays/internal/handlers/.


