from pyeoskit import wasmcompiler

def compile(name, code):
    return wasmcompiler.compile_go_src(name, code, replace=None)

def compile_file(name, go_file):
    with open(go_file, 'r') as f:
        code = f.read()
        return wasmcompiler.compile_go_src(name, code, replace=None)
    return None, None

def test_compiler():
    contract_files = [
        "testtransaction.go",
        "testpacksize.go",
        "testaction.go",
        "testmath.go",
        "testsingleton.go",
        "testasset.go",
        "testmi.go",
        "testsort.go",
        "testchain.go",
        "testtoken.go",
        "testcrypto.go",
        "testprimarykey.go",
        "testdb.go",
        "testprint.go",
        "testuint128.go",
        "testfloat128.go",
        "testprivileged.go",
        "testwasm.go",
        "testkv.go",
        "testlargecode.go",
        "testserializer.go"
    ]

    for file in contract_files:
        print(file)
        code, abi = compile_file(file[:-3], file)
        assert code
