module rcw

go 1.24.2

require github.com/rwinkhart/peercred-mini v0.0.0-20250407024456-ffe191394a3a

require (
	github.com/Microsoft/go-winio v0.6.2
	golang.org/x/sys v0.32.0
)

replace golang.org/x/sys => github.com/rwinkhart/sys-freebsd-13-xucred v0.0.0-20250405010723-99a5f0732c0e
