server:
	- go run cmd/main.go

test-coverage:
	mkdir -p repository/coverage
	cd repository && go test ./... -race -v -tags UnitTest -coverprofile ./coverage/coverage.out && go tool cover -func ./coverage/coverage.out && go tool cover -html ./coverage/coverage.out -o ./coverage/coverage.html
