tinygo build -x -gc=leaking -target eosio -wasm-abi=generic -scheduler=none -opt 0 -tags=math_big_pure_go -gen-code=true -strip=false -o test.wasm .
eosio-wasm2wast test.wasm -o test.wast
