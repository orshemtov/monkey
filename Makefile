build:
	go build -o out/monkey

test:
	go test

run:
	go run .

record: build
	vhs out/demo.tape
