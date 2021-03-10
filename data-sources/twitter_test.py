import unittest
import twitter


class TestTwitter(unittest.TestCase):

    def test_get_urls_from_tweet(self):
        tweet = '1368883070590476292'
        url = twitter.get_urls_from_tweet(tweet)
        self.assertEqual(['https://t.co/bLokglOAel'], url)

    def test_validate_json(self):
        jsons = [
            {
                'name': 'Valid JSON',
                'json': {
                    'address': 'desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu',
                    'pub_key': '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'signature': 'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    'value': 'ricmontagnin'
                },
                'valid': True
            },
            {
                'name': 'Missing address',
                'json': {
                    'pub_key': '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'signature': 'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    'value': 'ricmontagnin'
                },
                'valid': False
            },
            {
                'name': 'Missing pub_key',
                'json': {
                    'address': 'desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu',
                    'signature': 'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    'value': 'ricmontagnin'
                },
                'valid': False
            },
            {
                'name': 'Missing signature',
                'json': {
                    'address': 'desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu',
                    'pub_key': '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'value': 'ricmontagnin'
                },
                'valid': False
            },
            {
                'name': 'Missing value',
                'json': {
                    'address': 'desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu',
                    'pub_key': '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'signature': 'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                },
                'valid': False
            },
        ]

        for json in jsons:
            result = twitter.validate_json(json['json'])
            self.assertEqual(json['valid'], result, json['name'])

    def test_get_signature_from_url(self):
        # Valid signature
        result = twitter.get_signature_from_url('https://pastebin.com/raw/xz4S8WrW')
        self.assertEqual(True, result['valid'])

        # Invalid signature
        result = twitter.get_signature_from_url('https://bitcoin.org')
        self.assertEqual(False, result['valid'])

    def test_verify_signature(self):
        tests = [
            {
                'name': 'Valid data',
                'value': 'ricmontagnin',
                'signature': 'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                'pub_key': '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                'valid': True
            },
            {
                'name': 'Invalid value',
                'value': 'ricmontagnini',
                'signature': 'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                'pub_key': '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                'valid': False
            },
            {
                'name': 'Invalid signature',
                'value': 'ricmontagnin',
                'signature': 'a00a7d5bd45e5fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                'pub_key': '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                'valid': False
            },
            {
                'name': 'Invalid pub key',
                'value': 'ricmontagnin',
                'signature': 'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                'pub_key': '03302418cd13b082e7a7bc3ed05312a0b417d',
                'valid': False
            },
        ]

        for test in tests:
            result = twitter.verify_signature(test['pub_key'], test['signature'], test['value'])
            self.assertEqual(test['valid'], result, test['name'])


if __name__ == '__main__':
    unittest.main()
