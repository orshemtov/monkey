build:
	go build -o out/monkey

test:
	go test

run:
	go run .

record: build
	vhs out/demo.tape

record-all: build
	@for tape in out/demo-*.tape; do \
		echo "Recording $$tape..."; \
		vhs "$$tape"; \
	done
	@echo "All demos recorded!"

.PHONY: build test run record record-all
