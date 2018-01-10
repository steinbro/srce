export GOPATH=$(abspath ./bin/)
TARGETS=$(notdir $(wildcard ./bin/cmd/*))

all: $(TARGETS)

clean:
	rm exe/*

$(TARGETS):
	go build -o exe/$@ ./bin/cmd/$@

test:
	go test ./bin/src/srce/ -cover
