build/upload: 			*/*.go
	@go build -o build/upload ./cmd/main.go

run:					build/upload
	@build/upload

test-db: 					build/test-db
	@build/test-db -test.v

build/test-db: 			test/db_test.go
	@go test -c ./test -o build/test-db

clean:
	rm -rf build

.PHONY: test-db run


