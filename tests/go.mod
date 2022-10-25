module github.com/uuosio/chain/tests

go 1.17

replace github.com/uuosio/chain => ../

// replace github.com/uuosio/chaintester => ../../chaintester

require github.com/uuosio/chain v0.2.1

replace (
	github.com/uuosio/chain/tests => ./
	github.com/uuosio/chain/tests/testaction => ./testaction
	github.com/uuosio/chain/tests/testasset => ./testasset
	github.com/uuosio/chain/tests/testcrypto => ./testcrypto
	github.com/uuosio/chain/tests/testdb => ./testdb
	github.com/uuosio/chain/tests/testfloat128 => ./testfloat128
	github.com/uuosio/chain/tests/testlargecode => ./testlargecode
	github.com/uuosio/chain/tests/testmath => ./testmath
	github.com/uuosio/chain/tests/testmi => ./testmi
	github.com/uuosio/chain/tests/testpacksize => ./testpacksize
	github.com/uuosio/chain/tests/testprimarykey => ./testprimarykey
	github.com/uuosio/chain/tests/testprint => ./testprint
	github.com/uuosio/chain/tests/testprivileged => ./testprivileged
	github.com/uuosio/chain/tests/testsingleton => ./testsingleton
	github.com/uuosio/chain/tests/testsort => ./testsort
	github.com/uuosio/chain/tests/testtoken => ./testtoken
	github.com/uuosio/chain/tests/testtransaction => ./testtransaction
	github.com/uuosio/chain/tests/testuint128 => ./testuint128
	github.com/uuosio/chain/tests/testvariant => ./testvariant
)

require (
	github.com/apache/thrift v0.16.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/uuosio/chain/tests/testaction v0.0.1 // indirect
	github.com/uuosio/chain/tests/testasset v0.0.1 // indirect
	github.com/uuosio/chain/tests/testcrypto v0.0.1 // indirect
	github.com/uuosio/chain/tests/testdb v0.0.1 // indirect
	github.com/uuosio/chain/tests/testfloat128 v0.0.1 // indirect
	github.com/uuosio/chain/tests/testlargecode v0.0.1 // indirect
	github.com/uuosio/chain/tests/testmath v0.0.1 // indirect
	github.com/uuosio/chain/tests/testmi v0.0.1 // indirect
	github.com/uuosio/chain/tests/testpacksize v0.0.1 // indirect
	github.com/uuosio/chain/tests/testprimarykey v0.0.1 // indirect
	github.com/uuosio/chain/tests/testprint v0.0.1 // indirect
	github.com/uuosio/chain/tests/testprivileged v0.0.1 // indirect
	github.com/uuosio/chain/tests/testsingleton v0.0.1 // indirect
	github.com/uuosio/chain/tests/testsort v0.0.1 // indirect
	github.com/uuosio/chain/tests/testtoken v0.0.1 // indirect
	github.com/uuosio/chain/tests/testtransaction v0.0.1 // indirect
	github.com/uuosio/chain/tests/testuint128 v0.0.1 // indirect
	github.com/uuosio/chain/tests/testvariant v0.0.1 // indirect
	github.com/uuosio/chaintester v0.0.0-20221025012659-f8240d66cb56 // indirect
)
