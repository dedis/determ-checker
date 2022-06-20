# Golang static checker
A static checker that detects non-determinism in Golang smart contracts.

The static checker recursively parses all the source files of the smart contract, and issues a warning if the code uses any of the blacklisted types or if it imports any packages besides the whitelisted ones. We blacklisted _float_ types, _chan_ type, go routines _GoStmt_, and _map range_ statements. The list of whitelisted packages contains _fmt_, _bytes_, _hash_, _encoding_. 
