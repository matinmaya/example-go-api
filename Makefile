debug:
	GIN_MODE=debug go run cmd/app/main.go

test:
	GIN_MODE=test go run cmd/app/main.go

release:
	GIN_MODE=release go run cmd/app/main.go

seed:
	go run cmd/seeder/seeder.go