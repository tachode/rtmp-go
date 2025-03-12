all: test

test:
	go test -race -coverprofile=c.out ./...

clean:
	find . -name 'c.out' -exec rm -f '{}' \+
