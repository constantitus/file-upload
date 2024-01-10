build/upload: 			*/*.go
	@go build -o build/upload main/cmd


run:					build/upload
	@build/upload


dbcli: 				cmd/dbcli/*.go db/*.go
	@go build -o build/dbcli main/cmd/dbcli


db_test: 				build/db_test
	@build/test-db -test.v

build/db_test: 			test/db_test.go
	@go test -c main/test -o build/db_test

clean:
	rm -rf build

.PHONY: db_test run


