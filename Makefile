.PHONY: install-deps generate

install-deps:
	go install github.com/swaggo/swag/cmd/swag@latest

generate:
	go generate ./...
	#swag init --generalInfo ./cmd/http.go
	swag init -g ./cmd/http.go --parseDependency

swag:
	swag init -g ./cmd/http.go --parseDependency