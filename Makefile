.PHONY: lint
lint:
	cd backend && golangci-lint run