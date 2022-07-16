module test

go 1.17

replace github.com/uuosio/chain => /Users/newworld/dev/gscdk/chain

replace github.com/learnforpractice/chaintester => /Users/newworld/dev/gscdk/chaintester

require (
	github.com/learnforpractice/chaintester v0.0.0-20220711072919-951ecab8775b
	github.com/uuosio/chain v0.1.14
)

require (
	github.com/apache/thrift v0.16.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
)
