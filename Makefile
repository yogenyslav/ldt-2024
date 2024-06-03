.PHONY: lint
lint:
	cd api && golangci-lint run