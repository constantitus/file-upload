module main

go 1.21.4

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/go-pkgz/expirable-cache/v2 v2.0.0
	github.com/mattn/go-sqlite3 v1.14.18
	golang.org/x/time v0.5.0
	nullprogram.com/x/uuid v1.2.1
)

require nullprogram.com/x/isaac64 v1.0.0 // indirect

replace github.com/go-pkgz/expirable-cache/v2 => github.com/constantitus/expirable-cache/v2 v2.0.0-20231217195258-3fd0eed80adb
