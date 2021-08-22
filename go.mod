module github.com/ddirect/filetest

go 1.16

replace github.com/ddirect/xrand => ../xrand

replace github.com/ddirect/check => ../check

replace github.com/ddirect/format => ../format

require (
	github.com/ddirect/check v0.0.0-00010101000000-000000000000
	github.com/ddirect/format v0.0.0-00010101000000-000000000000
	github.com/ddirect/xrand v0.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.5.5
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
)
