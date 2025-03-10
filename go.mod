module verify_password

go 1.21.5
toolchain go1.23.7

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.36.0
)

require filippo.io/edwards25519 v1.1.0 // indirect
