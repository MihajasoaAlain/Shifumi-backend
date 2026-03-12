.PHONY: swag run
swag: 
			 swag init -g cmd/server/main.go --parseDependency --parseInternal   

run:  swag
			 go run cmd/server/main.go