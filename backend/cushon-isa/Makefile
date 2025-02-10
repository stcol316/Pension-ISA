build:
	@echo "Building..."
	@go build -o main cmd/investment-server/main.go

run: build
	@go run cmd/investment-server/main.go

test:
	@echo "Testing..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -f main