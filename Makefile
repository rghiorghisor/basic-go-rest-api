MODULE = $(shell go list -m)
PACKAGES := $(shell go list ./...)
GOLINT := $(shell go list -f {{.Target}} golang.org/x/lint/golint)

clean:
	rm -rf .bin cover.out cover-all.out

test:
	@echo "mode: count" > cover-all.out
	@$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=cover.out ${pkg}; \
		tail -n +2 cover.out >> cover-all.out;)

test-cover: test
	go tool cover -html=cover-all.out

lint:
	@$(GOLINT) ./...

build:
	go mod download && \
		go mod verify && \
		CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app cmd/api/main.go

run-dev:
	go run cmd/api/main.go

run-prod:
	go run cmd/api/main.go -env="production"