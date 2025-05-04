module github.com/rwinkhart/rcw

go 1.24.2

require github.com/rwinkhart/peercred-mini v0.0.0-20250407033241-c09add2eceea

require (
	github.com/Microsoft/go-winio v0.6.2
	golang.org/x/crypto v0.37.0
	golang.org/x/sys v0.32.0
)

replace golang.org/x/sys => github.com/rwinkhart/sys-freebsd-13-xucred v0.32.0

replace github.com/Microsoft/go-winio => github.com/rwinkhart/go-winio-easy-pipe-handles v0.0.0-20250407031321-96994a0e8410
