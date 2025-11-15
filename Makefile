build:
	go build -o out/monkey

test:
	go test

run:
	go run .

record: build
	vhs out/demo.tape
	cp out/demo.gif .github/demo.gif

.PHONY: build test run record 
