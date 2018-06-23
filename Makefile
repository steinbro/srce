export GOPATH=$(abspath ./bin/)
TARGETS=$(notdir $(wildcard ./bin/cmd/*))

all: $(TARGETS)

clean:
	rm exe/*

$(TARGETS):
	go build -o exe/$@ ./bin/cmd/$@

test:
	gofmt -w ./bin/
	go test ./bin/src/srce/ -cover

coverage:
	go test ./bin/src/srce/ -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out
