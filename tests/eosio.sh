# pushd /Users/newworld/Library/Caches/tinygo/goroot-go1.16.6-49c9e2ceb576a1ae8d6f12b11587ec6df46230fd2a2d481ba80455c3c4855816-syscall/src/eosiolib
# git pull
# popd

#rm -rf /Users/newworld/Library/Caches/tinygo/goroot-go1.16.6-49c9e2ceb576a1ae8d6f12b11587ec6df46230fd2a2d481ba80455c3c4855816-syscall/src/chain
#mkdir -p /Users/newworld/Library/Caches/tinygo/goroot-go1.16.6-49c9e2ceb576a1ae8d6f12b11587ec6df46230fd2a2d481ba80455c3c4855816-syscall/src/chain
#cp -r /Users/newworld/dev/github/go-eosiolib/* /Users/newworld/Library/Caches/tinygo/goroot-go1.16.6-49c9e2ceb576a1ae8d6f12b11587ec6df46230fd2a2d481ba80455c3c4855816-syscall/src/chain

tinygo build  -x -gc=leaking -o eosio.wasm -target eosio -wasm-abi=generic -scheduler=none  -opt z -tags=math_big_pure_go $1 $2 $3 $4
if [ $? -ne 0 ]; then
    echo "build failed"
    exit 1
fi
eosio-wasm2wast -o eosio.wast eosio.wasm
eosio-wast2wasm -o eosio2.wasm eosio.wast
