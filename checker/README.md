How to test the checker?

`go test checker_test.go -pkg=../inputs/whitelist-pkg.txt -types=../inputs/blacklist-types.txt -src=../inputs/test.go`

Whitelisted packages:
* fmt
* bytes
* hash
* encoding

Blacklisted types:
* goroutines & chan
* float
