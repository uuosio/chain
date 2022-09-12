
pushd $1
go-contract gencode -p $1 || exit 1
popd
