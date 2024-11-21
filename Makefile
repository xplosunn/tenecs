generate:
	go generate ./...

test:
	go test -count=1 ./... && npx cypress run