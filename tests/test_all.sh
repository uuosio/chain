run-ipyeos -m pytest -x -s test.py -k  test_hello || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_pack_size || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_action || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_crypto || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_mi || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_print || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_privileged || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_deffered_tx || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_db || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_token || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_singleton || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_asset || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_serializer || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_kv || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_primarykey || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_float128 || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_sort || exit 1

# run-ipyeos -m pytest -x -s test.py -k  test_largecode

run-ipyeos -m pytest -x -s test.py -k  test_math || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_go_math || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_malloc || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_uint128 || exit 1
run-ipyeos -m pytest -x -s test.py -k  test_revert || exit 1
