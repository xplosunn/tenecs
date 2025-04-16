build:
	go build -o tenecs cli/main.go

generate:
	go generate ./...

test:
	CI=true go test -count=1 ./... && npx cypress run

update_snaps:
	UPDATE_SNAPS=true go test ./...