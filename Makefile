all: build

build:
	go build cmd/number_requests/number_requests.go

test:
	go test -v -count=1 ./internal/...

clean:
	go clean
	rm -fr ./number_requests
