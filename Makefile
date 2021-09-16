.PHONY: clean test build_docker

clean:
	rm -f data/*
test:
	go test ./...
build_docker:
	docker build -t pensionera .