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
        self.chain.produce_block()

    def compile(cls, name, code):
        replace = '/Users/newworld/dev/github/go-chain'
        return wasmcompiler.compile_go_src(name, code, replace=replace)


    def test_hello(self):
        with open('test.wasm', 'rb') as f:
            code = f.read()
        with open('test.abi', 'r') as f:
            abi = json.load(f)
        self.chain.deploy_contract('hello', code, abi, 0)

        for i in range(4):
            r = self.chain.push_action('hello', 'sayhello', b'hello,world')
            print_console(r)
            self.chain.produce_block()
