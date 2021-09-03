import os
import sys
import json
import time
from inspect import currentframe, getframeinfo

from uuoskit import wasmcompiler

test_dir = os.path.dirname(__file__)
sys.path.append(os.path.join(test_dir, '..'))

from uuosio import log
from uuosio.chaintester import ChainTester

logger = log.get_logger(__name__)

def get_line_number():
    cf = currentframe()
    return cf.f_back.f_lineno

def print_console(tx):
    cf = currentframe()
    filename = getframeinfo(cf).filename

    num = cf.f_back.f_lineno

    if 'processed' in tx:
        tx = tx['processed']
    for trace in tx['action_traces']:
        # logger.info(trace['console'])
        print(f'+++++console:{num}', trace['console'])

        if not 'inline_traces' in trace:
            continue
        for inline_trace in trace['inline_traces']:
            # logger.info(inline_trace['console'])
            print(f'+++++console:{num}', inline_trace['console'])

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

    def compile(cls, name, code):
        replace = '/Users/newworld/dev/github/go-chain'
        return wasmcompiler.compile_go_src(name, code, replace=replace)

    def test_hello(self):
        code = '''
package main
import "chain/logger"
func main() {
    logger.Println("Hello,world!")
}
'''
        code, abi = self.compile('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_pack_size(self):
        with open('testpacksize.go', 'r') as f:
            code = f.read()
        code, abi = self.compile('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_action(self):
        with open('testaction.go', 'r') as f:
            code = f.read()
        code, abi = self.compile('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        try:
            r = self.chain.push_action('hello', 'sayhello', b'hello,world')
            print_console(r)
        except Exception as e:
            print_except(e.args[0])

        try:
            old_balance = self.chain.get_balance('hello')
            r = self.chain.push_action('hello', 'sayhello3', b'hello,world')
            print_console(r)
            new_balance = self.chain.get_balance('hello')
            assert abs(new_balance + 1.0 - old_balance) < 1e-9
        except Exception as e:
            print_except(e.args[0])

    def test_crypto(self):
        with open('testcrypto.go', 'r') as f:
            code = f.read()
        code, abi = self.compile('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, abi, 0)

        r = self.chain.push_action('hello', 'testhash', '')
        self.chain.produce_block()

        sig = 'SIG_K1_KiXXExwMGG5NvAngS3X58fXVVcnmPc7fxgwLQAbbkSDj9gwcxWHxHwgpUegSCfgp4nFMMgjLDAKSQWZ2NLEmcJJn1m2UUg'
        pub = 'EOS7wy4M8ZTYqtoghhDRtE37yRoSNGc6zC2zFgdVmaQnKV5ZXe4kV'
        data = b'hello,world'
        args = {
            'data': data.hex(),
            'sig': sig,
            'pub': pub,
        }
        r = self.chain.push_action('hello', 'testrecover', args)
        print_console(r)

    def test_mi(self):
        with open('testmi.go', 'r') as f:
            code = f.read()
        code, abi = self.compile('testmi', code)
        logger.info("++++++++++code size %f", len(code)/1024)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'test1', b'hello,world')
        print_console(r)

        self.chain.produce_block()
        r = self.chain.push_action('hello', 'test2', b'hello,world')
        print_console(r)

    def test_print(self):
        with open('testprint.go', 'r') as f:
            code = f.read()
        code, abi = self.compile('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_privileged(self):
        with open('testprivileged.go', 'r') as f:
            code = f.read()
        code, abi = self.compile('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_deffered_tx(self):
        with open('testtransaction.go', 'r') as f:
            code = f.read()
        code, abi = self.compile('testtransaction', code)
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
        code, abi = self.compile('testdb', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_token(self):
        with open('testtoken.go', 'r') as f:
            code = f.read()
        code, abi = self.compile('testtoken', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_singleton(self):
        with open('testsingleton.go', 'r') as f:
            code = f.read()
        code, abi = self.compile('testsingleton', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)

        for i in range(4):
            r = self.chain.push_action('hello', 'sayhello', b'hello,world')
            print_console(r)
            self.chain.produce_block()

    def test_asset(self):
        with open('testasset.go', 'r') as f:
            code = f.read()
        code, abi = self.compile('testasset', code)
        assert code
        self.chain.deploy_contract('hello', code, b'', 0)

        try:
            r = self.chain.push_action('hello', 'test1', b'hello,world')
        except Exception as e:
            error_msg = e.args[0]['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'addition overflow'

        try:
            r = self.chain.push_action('hello', 'test2', b'hello,world')
        except Exception as e:
            error_msg = e.args[0]['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'subtraction underflow'
        self.chain.produce_block()

        # magnitude of asset amount must be less than 2^62
        try:
            r = self.chain.push_action('hello', 'test3', b'hello,world')
        except Exception as e:
            error_msg = e.args[0]['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'magnitude of asset amount must be less than 2^62'
        self.chain.produce_block()

        #divide by zero
        try:
            r = self.chain.push_action('hello', 'test4', b'hello,world')
        except Exception as e:
            error_msg = e.args[0]['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'divide by zero'
        self.chain.produce_block()

        #signed division overflow
        try:
            r = self.chain.push_action('hello', 'test5', b'hello,world')
        except Exception as e:
            error_msg = e.args[0]['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'signed division overflow'
        self.chain.produce_block()

        #bad symbol
        try:
            r = self.chain.push_action('hello', 'test11', b'hello,world')
        except Exception as e:
            error_msg = e.args[0]['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'bad symbol'
        self.chain.produce_block()

        #multiplication overflow
        try:
            r = self.chain.push_action('hello', 'test12', b'hello,world')
        except Exception as e:
            error_msg = e.args[0]['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'multiplication overflow'
        self.chain.produce_block()

        #multiplication underflow
        try:
            r = self.chain.push_action('hello', 'test13', b'hello,world')
        except Exception as e:
            error_msg = e.args[0]['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'multiplication underflow'
        self.chain.produce_block()
