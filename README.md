# deter-checker
A static checker for non-determinism in source files

go test checker_test.go -pkg=../inputs/whitelist-pkg.txt -types=../inputs/blacklist-types.txt -src=../inputs/test.go
