all: test

install-swagger:
	brew install swagger-codegen@2
test: get-deps generate
	go test ./...

get-deps:
	go get -u golang.org/x/time/rate
	go get -u github.com/stretchr/testify
	go get -u gopkg.in/h2non/gock.v1

generate:
	go generate
