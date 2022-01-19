import unittest

import httpretty

import github


class GitHubTest(unittest.TestCase):

    @httpretty.activate(verbose=True, allow_net_connect=False)
    def test_get_data_from_gist(self):
        # Register fake HTTP call
        httpretty.register_uri(
            httpretty.GET,
            "https://gist.githubusercontent.com/RiccardoM/720e0072390a901bb80e59fd60d7fded/raw/",
            status=200,
            body='{"address":"8902A4822B87C1ADED60AE947044E614BD4CAEE2","pub_key":"033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d","value":"Riccardo Montagnin#5414","signature":"d10db146bb4d234c5c1d2bc088e045f4f05837c690bce4101e2c0f0c6c96e1232d8516884b0a694ee85e9c9da51be74966886cbb12af4ad87e5336da76d75cfb"}',
        )

        # Valid signature
        data = github.get_data_from_gist(github.CallData('RiccardoM', '720e0072390a901bb80e59fd60d7fded'))
        self.assertIsNotNone(data)

        # Register fake HTTP call
        httpretty.register_uri(
            httpretty.GET,
            "https://gist.githubusercontent.com/RiccardoM/invalid_gist_id/raw/",
            status=200,
            body='{}',
        )

        # Invalid signature
        data = github.get_data_from_gist(github.CallData('RiccardoM', 'invalid_gist_id'))
        self.assertIsNone(data)

    def test_validate_json(self):
        jsons = [
            {
                'name': 'Valid JSON',
                'json': {
                    'address': '8902A4822B87C1ADED60AE947044E614BD4CAEE2',
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
                    'address': '8902A4822B87C1ADED60AE947044E614BD4CAEE2',
                    'signature': 'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    'value': 'ricmontagnin'
                },
                'valid': False
            },
            {
                'name': 'Missing signature',
                'json': {
                    'address': '8902A4822B87C1ADED60AE947044E614BD4CAEE2',
                    'pub_key': '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'value': 'ricmontagnin'
                },
                'valid': False
            },
            {
                'name': 'Missing value',
                'json': {
                    'address': '8902A4822B87C1ADED60AE947044E614BD4CAEE2',
                    'pub_key': '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'signature': 'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                },
                'valid': False
            },
        ]

        for json in jsons:
            result = github.validate_json(json['json'])
            self.assertEqual(json['valid'], result, json['name'])

    def test_verify_signature(self):
        tests = [
            {
                'name': 'Valid data',
                'valid': True,
                'data': github.VerificationData(
                    '',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    '7269636d6f6e7461676e696e',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                ),
            },
            {
                'name': 'Invalid value',
                'valid': False,
                'data': github.VerificationData(
                    '',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    'ricmontagni',
                ),
            },
            {
                'name': 'Invalid signature',
                'valid': False,
                'data': github.VerificationData(
                    '',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'a00a7d5bd45e42615645fcaeb4d800af2704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    '7269636d6f6e7461676e696e',
                ),
            },
            {
                'name': 'Invalid pub key',
                'valid': False,
                'data': github.VerificationData(
                    '',
                    '033024e9e0ad4f9305ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    '7269636d6f6e7461676e696e',
                ),
            },
        ]

        for test in tests:
            result = github.verify_signature(test['data'])
            self.assertEqual(test['valid'], result, test['name'])

    def test_verify_address(self):
        tests = [
            {
                'name': 'Valid address',
                'valid': True,
                'data': github.VerificationData(
                    '8902A4822B87C1ADED60AE947044E614BD4CAEE2',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    '7269636d6f6e7461676e696e',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173'
                ),
            },
            {
                'name': 'Invalid address',
                'valid': False,
                'data': github.VerificationData(
                    '8902A4822B87C1ADED60AE947044E614BD4CAEE2',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b41',
                    '7269636d6f6e7461676e696e',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173'
                ),
            },
        ]

        for test in tests:
            result = github.verify_address(test['data'])
            self.assertEqual(test['valid'], result, test['name'])


if __name__ == '__main__':
    unittest.main()
