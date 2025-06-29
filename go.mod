module github.com/rwinkhart/rcw

go 1.24.4

require (
	github.com/Microsoft/go-winio v0.6.2
	github.com/rwinkhart/go-boilerplate v0.1.0
	github.com/rwinkhart/peercred-mini v0.1.0
	golang.org/x/crypto v0.39.0
	golang.org/x/sys v0.33.0
)

require golang.org/x/term v0.32.0 // indirect

replace golang.org/x/sys => github.com/rwinkhart/sys-freebsd-13-xucred v0.33.0

replace github.com/Microsoft/go-winio => github.com/rwinkhart/go-winio-easy-pipe-handles v0.0.0-20250407031321-96994a0e8410
