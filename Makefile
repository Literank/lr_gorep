.PHONY: run test clean

run:
	go run main.go -r -n result .

test:
	go test -v ./...

install:
	go install

build:
	go build -o gorep

clean:
	rm -rf gorep
