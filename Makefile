up:
	docker compose -f ./deployments/docker-compose.yml up -d --build

down:
	docker compose -f ./deployments/docker-compose.yml down

make swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --parseDependency --parseInternal -g cmd/main.go

test:
	go test credit_holidays/...

wire:
	go install github.com/google/wire/cmd/wire@latest
	wire credit_holidays/internal/handlers/.

generate-mocks:
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen
	mockgen -destination=internal/mocks/mock_db.go -package=mocks credit_holidays/internal/db CreditHolidaysDB
