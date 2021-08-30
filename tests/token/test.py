import os
import sys
import json
import time

from uuoskit import wasmcompiler, uuosapi

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

    def test_token(self):
        with open('token.wasm', 'rb') as f:
            code = f.read()
        with open('token.abi', 'r') as f:
            abi = f.read()
        self.chain.deploy_contract('hello', code, abi, 0)
        create = {
            "issuer": "hello",
            "maximum_supply": "100.0000 EEOS",
        }
        r = self.chain.push_action('hello', 'create', create)
        print_console(r)
        logger.info('+++++++create elapsed: %s', r['elapsed'])
        self.chain.produce_block()

        r = self.chain.get_table_rows(True, 'hello', 'EEOS', 'stat', "", "")
        logger.info(r)
        assert r['rows'][0]['Issuer'] == 'hello'
        assert r['rows'][0]['MaxSupply'] == '100.0000 EEOS'
        assert r['rows'][0]['Supply'] == '0.0000 EEOS'

        try:
            r = self.chain.push_action('hello', 'create', create)
        except Exception as e:
            error_msg = e.args[0]['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'token with symbol already exists'
            # logger.info(json.dumps(e.args[0], indent='    '))

        #test issue

        issue = {'to': 'hello', 'quantity': '1.0000 EEOS', 'memo': 'issue to alice'}
        r = self.chain.push_action('hello', 'issue', issue)
        logger.info('+++++++issue elapsed: %s', r['elapsed'])
        self.chain.produce_block()

        r = self.chain.get_table_rows(True, 'hello', 'EEOS', 'stat', "", "")
        logger.info(r)
        assert r['rows'][0]['Issuer'] == 'hello'
        assert r['rows'][0]['MaxSupply'] == '100.0000 EEOS'
        assert r['rows'][0]['Supply'] == '1.0000 EEOS'

        r = self.chain.get_table_rows(True, 'hello', 'hello', 'accounts', "", "")
        logger.info(r)
        assert r['rows'][0]['Balance'] == '1.0000 EEOS'

        try:
            issue = {'to': 'eosio', 'quantity': '1.0000 EEOS', 'memo': 'issue to alice'}
            self.chain.push_action('hello', 'issue', issue)
        except Exception as e:
            error_msg = e.args[0]['action_traces'][0]['except']['stack'][0]['data']['s']
            assert error_msg == 'tokens can only be issued to issuer account'

        #test transfer
        transfer = {'from': 'hello', 'to': 'alice', 'quantity': '1.0000 EEOS', 'memo': 'transfer from alice'}
        r = self.chain.push_action('hello', 'transfer', transfer)
        logger.info('+++++++transfer elapsed: %s', r['elapsed'])

        self.chain.produce_block()

        r = self.chain.get_table_rows(True, 'hello', 'hello', 'accounts', "", "")
        logger.info(r)
        assert r['rows'][0]['Balance'] == '0.0000 EEOS'

        r = self.chain.get_table_rows(True, 'hello', 'alice', 'accounts', "", "")
        logger.info(r)
        assert r['rows'][0]['Balance'] == '1.0000 EEOS'

        # transfer back
        transfer = {'from': 'alice', 'to': 'hello', 'quantity': '1.0000 EEOS', 'memo': 'transfer back'}
        r = self.chain.push_action('hello', 'transfer', transfer, {'alice': 'active'})
        logger.info('+++++++transfer elapsed: %s', r['elapsed'])
        self.chain.produce_block()

        #quantity chain.Asset, memo
        retire = {'quantity': '1.0000 EEOS', 'memo': 'retire 1.0000 EEOS'}
        r = self.chain.push_action('hello', 'retire', retire)
        logger.info('+++++++retire elapsed: %s', r['elapsed'])

        r = self.chain.get_table_rows(True, 'hello', 'hello', 'accounts', "", "")
        assert r['rows'][0]['Balance'] == '0.0000 EEOS'

        r = self.chain.get_table_rows(True, 'hello', 'EEOS', 'stat', "", "")
        logger.info(r)
        assert r['rows'][0]['Supply'] == '0.0000 EEOS'


        r = self.chain.get_table_rows(True, 'hello', 'helloworld11', 'accounts', "", "")
        assert len(r['rows']) == 0

        #owner chain.Name, symbol chain.Symbol, ram_payer chain.Name
        #test open
        open_action = {'owner': 'helloworld11', 'symbol': '4,EEOS', 'ram_payer': 'hello'}
        r = self.chain.push_action('hello', 'open', open_action)
        logger.info('+++++++open elapsed: %s', r['elapsed'])

        r = self.chain.get_table_rows(True, 'hello', 'helloworld11', 'accounts', "", "")
        assert r['rows'][0]['Balance'] == '0.0000 EEOS'

        #test close
        close_action = {'owner': 'helloworld11', 'symbol': '4,EEOS'}
        r = self.chain.push_action('hello', 'close', close_action, {'helloworld11': 'active'})
        logger.info('+++++++close elapsed: %s', r['elapsed'])
        self.chain.produce_block()

        r = self.chain.get_table_rows(True, 'hello', 'helloworld11', 'accounts', "", "")
        assert len(r['rows']) == 0
