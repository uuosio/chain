import os
import sys
import json
import time

from uuoskit import wasmcompiler

test_dir = os.path.dirname(__file__)
sys.path.append(os.path.join(test_dir, '..'))

from uuosio import log
from uuosio.chaintester import ChainTester

logger = log.get_logger(__name__)

def print_console(tx):
    if 'processed' in tx:
        tx = tx['processed']
    for trace in tx['action_traces']:
        logger.info(trace['console'])
        if not 'inline_traces' in trace:
            continue
        for inline_trace in trace['inline_traces']:
            logger.info('++inline console:', inline_trace['console'])

def print_except(tx):
    if 'processed' in tx:
        tx = tx['processed']
    for trace in tx['action_traces']:
        logger.info(trace['console'])
        logger.info(json.dumps(trace['except'], indent=4))

class Test(object):

    @classmethod
    def setup_class(cls):
        cls.main_token = 'UUOS'
        cls.chain = ChainTester()

        test_account1 = 'hello'
        a = {
            "account": test_account1,
            "permission": "active",
            "parent": "owner",
            "auth": {
                "threshold": 1,
                "keys": [
                    {
                        "key": 'EOS6AjF6hvF7GSuSd4sCgfPKq5uWaXvGM2aQtEUCwmEHygQaqxBSV',
                        "weight": 1
                    }
                ],
                "accounts": [{"permission":{"actor":test_account1,"permission": 'eosio.code'}, "weight":1}],
                "waits": []
            }
        }
        cls.chain.push_action('eosio', 'updateauth', a, {test_account1:'active'})
        cls.chain.push_action('eosio', 'setpriv', {'account':'hello', 'is_priv': True}, {'eosio':'active'})

    @classmethod
    def teardown_class(cls):
        cls.chain.free()

    def setup_method(self, method):
        pass

    def teardown_method(self, method):
        pass

    def test_hello(self):
        code = '''
package main
import "chain/logger"
func main() {
    logger.Println("Hello,world!")
}
'''
        code, abi = wasmcompiler.compile_go_src('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_pack_size(self):
        with open('testpacksize.go', 'r') as f:
            code = f.read()
        code, abi = wasmcompiler.compile_go_src('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_action(self):
        with open('testaction.go', 'r') as f:
            code = f.read()
        code, abi = wasmcompiler.compile_go_src('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        try:
            r = self.chain.push_action('hello', 'sayhello', b'hello,world')
            print_console(r)
        except Exception as e:
            print_except(e.args[0])
            # logger.info(json.dumps(e.args[0], indent=4))
            # error = e.args[0]['except']
            # logger.info('error:', error)

    def test_crypto(self):
        with open('testcrypto.go', 'r') as f:
            code = f.read()
        code, abi = wasmcompiler.compile_go_src('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_mi(self):
        with open('testmi.go', 'r') as f:
            code = f.read()
        code, abi = wasmcompiler.compile_go_src('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_print(self):
        with open('testprint.go', 'r') as f:
            code = f.read()
        code, abi = wasmcompiler.compile_go_src('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_privileged(self):
        with open('testprivileged.go', 'r') as f:
            code = f.read()
        code, abi = wasmcompiler.compile_go_src('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_deffered_tx(self):
        with open('testtransaction.go', 'r') as f:
            code = f.read()
        code, abi = wasmcompiler.compile_go_src('testtransaction', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello1', b'hello,world')
        print_console(r)

        self.chain.produce_block()
        self.chain.produce_block()

#        time.sleep(1)
        r = self.chain.push_action('hello', 'sayhello3', b'hello,world')
        print_console(r)

    def test_db(self):
        with open('testdb.go', 'r') as f:
            code = f.read()
        code, abi = wasmcompiler.compile_go_src('testdb', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_token(self):
        with open('testtoken.go', 'r') as f:
            code = f.read()
        code, abi = wasmcompiler.compile_go_src('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)