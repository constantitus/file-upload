build/upload: 			*/*.go
	@go build -o build/upload main/cmd


run:					build/upload
	@build/upload


dbcli: 				cmd/dbcli/*.go db/*.go
	@go build -o build/dbcli main/cmd/dbcli


test-db: 				build/test-db
	@build/test-db -test.v

build/test-db: 			test/db_test.go
	@go test -c ./test -o build/test-db

clean:
	rm -rf build

.PHONY: test-db run


