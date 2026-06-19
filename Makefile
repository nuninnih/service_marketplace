# run app
app-run:
	go run app/api/main.go

# mock
mock-install:
	go install github.com/golang/mock/mockgen@latest
mock-user:
	mockgen -source service/user/user_repo.go -destination service/user/mock/user_mock_repo.go

# test
test-service:
	go test -v ./service/user/... -coverprofile=coverage.out -cover -failfast
test-service-coverage:
	go test -v $$(go list ./service/user/... | grep -v '/mock') -coverprofile=coverage.out -cover -failfast && \
	go tool cover -html=coverage.out -o cover.html && \
	open cover.html
