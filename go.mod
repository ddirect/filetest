module github.com/ddirect/filetest

go 1.16

replace github.com/ddirect/xrand => ../xrand

replace github.com/ddirect/check => ../check

require (
	github.com/ddirect/check v0.0.0-00010101000000-000000000000
	github.com/ddirect/xrand v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
)
