import os
import sys
import json
import time
import pytest
from inspect import currentframe, getframeinfo

from pyeoskit import wasmcompiler

test_dir = os.path.dirname(__file__)
sys.path.append(os.path.join(test_dir, '..'))

from ipyeos import log
from ipyeos.chaintester import ChainTester

from ipyeos import chaintester
from ipyeos.chain_exceptions import ChainException

chaintester.chain_config['contracts_console'] = True

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
        print(f'+++++console:{num}\n%s'%(trace['console'], ))

        if not 'inline_traces' in trace:
            continue
        for inline_trace in trace['inline_traces']:
            # logger.info(inline_trace['console'])
            print(f'+++++console:{num}\n%s'%(inline_trace['console'], ))

def print_except(tx):
    if 'processed' in tx:
        tx = tx['processed']
    for trace in tx['action_traces']:
        logger.info(trace['console'])
        logger.info(json.dumps(trace['except'], indent=4))

class MyChainTester(ChainTester):
    def push_action(self, account, action, args,  permissions=None, explicit_cpu_bill=False):
        if not permissions:
            permissions = {account: 'active'}
        return self.push_action_ex(account, action, args, permissions, explicit_cpu_bill)

def init_chain():
    chain = MyChainTester()
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
    chain.push_action('eosio', 'updateauth', a, {test_account1:'active'})
    chain.push_action('eosio', 'setpriv', {'account':'hello', 'is_priv': True}, {'eosio':'active'})
    return chain

class Test(object):

    @classmethod
    def setup_class(cls):
        pass
    @classmethod
    def teardown_class(cls):
        pass

    def setup_method(self, method):
        self.chain = init_chain()

    def teardown_method(self, method):
        self.chain.free()

    def compile(cls, name, code):
        replace = os.path.join(test_dir, '..')
        return wasmcompiler.compile_go_src(name, code, replace=replace)

    def test_hello(self):
        return
        code = '''
package main
import "github.com/uuosio/chain/logger"
func main() {
    logger.Println("Hello,world!")
}
'''
        code, abi = self.compile('hello', code)
        assert code
        self.chain.deploy_contract('hello', code, '', 0)
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_pack_size(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        self.chain.deploy_contract('hello', code, '', 0)
        self.chain.push_action("hello", "settest", b'testserializer')

        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_action(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testaction')

        try:
            r = self.chain.push_action('hello', 'sayhello', b'hello,world')
            print_console(r)
        except ChainException as e:
            print_except(e.json())

        try:
            old_balance = self.chain.get_balance('hello')
            r = self.chain.push_action('hello', 'sayhello3', b'hello,world')
            print_console(r)
            new_balance = self.chain.get_balance('hello')
            assert new_balance + 10000 == old_balance
        except ChainException as e:
            print_except(e.json())

    def test_crypto(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        with open('testcrypto/test.abi', 'rb') as f:
            abi = f.read()

        self.chain.deploy_contract('hello', code, abi)
        self.chain.push_action("hello", "settest", b'testcrypto')

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
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        with open('testmi/testmi.abi', 'rb') as f:
            abi = f.read()

        self.chain.deploy_contract('hello', code, abi)
        self.chain.push_action("hello", "settest", b'testmi')

        r = self.chain.push_action('hello', 'test1', b'hello,world')
        print_console(r)

        # self.chain.produce_block()
        # r = self.chain.push_action('hello', 'test2', b'hello,world')
        # print_console(r)

    def test_print(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testprint')

        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_privileged(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testprivileged')

        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_deffered_tx(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()
        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testtransaction')

        r = self.chain.push_action('hello', 'sayhello1', b'hello,world')
        print_console(r)

        self.chain.produce_block()
        self.chain.produce_block()
        self.chain.produce_block()

#        time.sleep(1)
        r = self.chain.push_action('hello', 'sayhello3', b'hello,world')
        print_console(r)

    def test_db(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()
        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testdb')
        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_token(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()
        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testtoken')

        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_singleton(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()
        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testsingleton')

        # for i in range(4):
        #     r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        #     print_console(r)
        #     self.chain.produce_block()

    def test_asset(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()
        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testasset')

        try:
            r = self.chain.push_action('hello', 'test1', b'hello,world')
        except ChainException as e:
            error_msg = e.json()['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'addition overflow'

        try:
            r = self.chain.push_action('hello', 'test2', b'hello,world')
        except ChainException as e:
            error_msg = e.json()['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'subtraction underflow'
        self.chain.produce_block()

        # magnitude of asset amount must be less than 2^62
        try:
            r = self.chain.push_action('hello', 'test3', b'hello,world')
        except ChainException as e:
            error_msg = e.json()['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'magnitude of asset amount must be less than 2^62'
        self.chain.produce_block()

        #divide by zero
        try:
            r = self.chain.push_action('hello', 'test4', b'hello,world')
        except Exception as e:
            error_msg = e.json()['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'divide by zero'
        self.chain.produce_block()

        #divide by negative value
        try:
            r = self.chain.push_action('hello', 'test5', b'hello,world')
        except Exception as e:
            error_msg = e.json()['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'divide by negative value'
        self.chain.produce_block()

        #bad symbol
        try:
            r = self.chain.push_action('hello', 'test11', b'hello,world')
        except Exception as e:
            error_msg = e.json()['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'bad symbol'
        self.chain.produce_block()

        #multiplication overflow
        try:
            r = self.chain.push_action('hello', 'test12', b'hello,world')
        except Exception as e:
            error_msg = e.json()['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'multiplication overflow'
        self.chain.produce_block()

        #multiplication underflow
        try:
            r = self.chain.push_action('hello', 'test13', b'hello,world')
        except Exception as e:
            error_msg = e.json()['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'multiplication underflow'
        self.chain.produce_block()

    def test_serializer(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        with open('testserializer/test.abi', 'r') as f:
            abi = f.read()

        self.chain.deploy_contract('hello', code, abi, 0)
        self.chain.push_action("hello", "settest", b'testserializer')
        
        args = dict(
            a0 = True, # a0 bool,
            a1 = 0xff, # a1 int8,
            a2 = 0xff, # a2 uint8,
            a3 = 0xffff, # a3 int16,
            a4 = 0xffff, # a4 uint16,
            a5 = 0xffffffff, # a5 int32,
            a6 = 0xffffffff, # a6 uint32,
            a7 = 0xffffffffffffffff, # a7 int64,
            a8 = 0xffffffffffffffff, # a8 uint64,
            a9 = '0x7fffffffffffffffffffffffffffffff', # // a9 int128,
            a10 = '0xffffffffffffffffffffffffffffffff', # a10 chain.Uint128,
            a11 = 0xffffffff, # // a11 varint32,
            a12 = 0xffffffff, # // a12 varuint32,
            a13 = 11.2233, # a13 float32,
            a14 = 11.2233, # a14 float64,
        	a15 = '0x7fffffffffffffffffffffffffffffff', #  a15 chain.Float128,
	        a16 = '2021-09-03T04:13:21', #  a16 chain.TimePoint,
        	a17 = '2021-09-03T04:13:21', # a17 chain.TimePointSec,
            a18 = '2021-09-03T04:13:21', # {'slot': 193723200}, #a18 chain.BlockTimeStamp, //block_timestamp_type,
            a19 = 'helloworld', # a19 chain.Name,
            a20 = b'hello,world'.hex(), # a20 []byte, //bytes,
            a21 = 'hello,world', # a21 string,
            a22 = 'aa'*20, # a22 chain.Checksum160, //checksum160,
            a23 = 'aa'*32, # a23 chain.Checksum256, //checksum256,
            a24 = 'aa'*64, # a24 chain.Checksum512, //checksum512,
            a25 = 'EOS5HoPaVaPivnVHsCvpoKZMmB6gcWGV5b3vF7S6pfsgFACzufMDy', # a25 chain.PublicKey, //public_key,
            a26 = 'SIG_K1_KbSF8BCNVA95KzR1qLmdn4VnxRoLVFQ1fZ8VV5gVdW1hLfGBdcwEc93hF7FBkWZip1tq2Ps27UZxceaR3hYwAjKL7j59q8', # a26 chain.Signature, //signature,
            a27 = '4,EOS',# a27 chain.Symbol, //symbol,
            a28 = 'EOS', # a28 chain.SymbolCode, //symbol_code,
            a29 = '1.0000 EOS', # a29 chain.Asset,
            a30 = {'quantity': '1.0000 EOS', 'contract': 'eosio.token'}
        )
        r = self.chain.push_action('hello', 'test', args)

        args = dict(
            a1 = -1,
            a2 = 0x7fffffff
        )
        r = self.chain.push_action('hello', 'testvarint', args)
        print_console(r)

        r = self.chain.push_action('hello', 'testpack', '')
        print_console(r)

        args = dict(
            a1 = 0xffffffff
        )
        r = self.chain.push_action('hello', 'testvaruint', args)
        print_console(r)

        try:
            r = self.chain.push_action('hello', 'testvaruint', b'')
            raise Exception("shoud not go here")
        except Exception as e:
            e = e.json()
            assert e['action_traces'][0]['except']['stack'][0]['data']['s'] == 'raw VarUint32 value can not be empty'

    def test_primarykey(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        with open('testprimarykey/test.abi', 'rb') as f:
            abi = f.read()

        self.chain.deploy_contract('hello', code, abi)
        self.chain.push_action("hello", "settest", b'testprimarykey')

        self.chain.push_action('hello', 'sayhello', b'')
        self.chain.produce_block()

        try:
            self.chain.push_action('hello', 'sayhello', b'')
        except Exception as e:
            error_msg = e.json()['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'mi.Update: Can not change primary key duration update'
        self.chain.produce_block()

    def test_float128(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testfloat128')

        r = self.chain.push_action('hello', 'sayhello', b'hello,world')
        print_console(r)

    def test_sort(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()
        with open('testsort/test.abi', 'rb') as f:
            abi = f.read()

        self.chain.deploy_contract('hello', code, abi)
        self.chain.push_action("hello", "settest", b'testsort')

        pubs = [
            "EOS6SD6yzqaZhdPHw2LUVmZxWLeWxnp76KLnnBbqP94TsDsjNLosG",
            "EOS4vtCi4jbaVCLVJ9Moenu9j7caHeoNSWgWY65bJgEW8MupWsRMo",
            "EOS82JTja1SbcUjSUCK8SNLLMcMPF8W5fwUYRXmX32obtjsZMW9nx"
        ]
        hex_pubs = []
        import base58
        for pub in pubs:
            h = base58.b58decode(pub[3:])[:-4].hex()
            hex_pubs.append(h)
            logger.info(h)

        r = self.chain.push_action('hello', 'test', {'pubs': pubs})
        print_console(r)


        logger.info(hex_pubs)
        hex_pubs.sort()
        logger.info(hex_pubs)

    @pytest.mark.skip(reason="deprecated")
    def test_largecode(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()
        with open('testlargecode/test.abi', 'rb') as f:
            abi = f.read()

        self.chain.deploy_contract('hello', code, abi)
        self.chain.push_action("hello", "settest", b'testsort')

        self.chain.deploy_contract('hello', code, abi, 0)
        r = self.chain.push_action('hello', 'test', b'hello,world')
        print_console(r)

    def test_go_math(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testmath')

        with pytest.raises(Exception) as e:
            r = self.chain.push_action('hello', 'test', b'hello,world')
            print_console(r)

    def test_uint128(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()

        self.chain.deploy_contract('hello', code, '')
        self.chain.push_action("hello", "settest", b'testuint128')

        r = self.chain.push_action('hello', 'test', b'hello,world')
        print_console(r)

    def test_variant(self):
        with open('tests.wasm', 'rb') as f:
            code = f.read()
        with open('testvariant/test.abi', 'rb') as f:
            abi = f.read()

        self.chain.deploy_contract('hello', code, abi)
        self.chain.push_action("hello", "settest", b'testvariant')

        args = {
            "v": ['uint64', 123]
        }
        r = self.chain.push_action('hello', 'test', args)
        print_console(r)
        ret = self.chain.get_table_rows(True, 'hello', '', 'mytable', '', '', 10)
        logger.info("%s", ret)
        assert ret['rows'][0]['a'] == ['uint64', 123]
