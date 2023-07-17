import unittest
import instagram
import httpretty


class instagramTest(unittest.TestCase):

    @httpretty.activate(verbose=True, allow_net_connect=False)
    def test_get_user_data(self):
        # Register fake HTTP call
        httpretty.register_uri(
            httpretty.GET,
            "https://themis.mainnet.desmos.network/instagram/test_user_desmos",
            status=200,
            body='{"address":"71b0310267b49279116835ed35791c24c110012f","pub_key":"0203233fabd69a1b7a90bb968a0ab66e3af61989f65cf0bc1f8e9518740a302f1f","value":"746573745f757365725f6465736d6f73","signature":"c12605456b8652df655bb43d0166586dfc0c5d758b03f127ca6b027d0ec140ca29b9569a20c9b78b72e13d15c1a7fa0b142dc0e624f3f51ef76bd94e55345d2a"}',
        )

        # Valid signature
        data = instagram.get_user_data(instagram.CallData('test_user_desmos'))
        self.assertIsNotNone(data)

        # Invalid signature
        httpretty.register_uri(
            httpretty.GET,
            "https://themis.mainnet.desmos.network/instagram/test_user_desmos",
            status=404,
        )
        data = instagram.get_user_data(instagram.CallData('test_user_desmos'))
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
            result = instagram.validate_json(json['json'])
            self.assertEqual(json['valid'], result, json['name'])

    def test_verify_signature(self):
        tests = [
            {
                'name': 'Valid data',
                'valid': True,
                'data': instagram.VerificationData(
                    '',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    '7269636d6f6e7461676e696e',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                ),
            },
            {
                'name': 'Invalid value',
                'valid': False,
                'data': instagram.VerificationData(
                    '',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    'ricmontagni',
                ),
            },
            {
                'name': 'Invalid signature',
                'valid': False,
                'data': instagram.VerificationData(
                    '',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'a00a7d5bd45e42615645fcaeb4d800af2704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    '7269636d6f6e7461676e696e',
                ),
            },
            {
                'name': 'Invalid pub key',
                'valid': False,
                'data': instagram.VerificationData(
                    '',
                    '033024e9e0ad4f9305ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    '7269636d6f6e7461676e696e',
                ),
            },
        ]

        for test in tests:
            result = instagram.verify_signature(test['data'])
            self.assertEqual(test['valid'], result, test['name'])

    def test_verify_address(self):
        tests = [
            {
                'name': 'Valid address',
                'valid': True,
                'data': instagram.VerificationData(
                    '8902A4822B87C1ADED60AE947044E614BD4CAEE2',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    '7269636d6f6e7461676e696e',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173'
                ),
            },
            {
                'name': 'Invalid address',
                'valid': False,
                'data': instagram.VerificationData(
                    '8902A4822B87C1ADED60AE947044E614BD4CAEE2',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b41',
                    '7269636d6f6e7461676e696e',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173'
                ),
            },
        ]

        for test in tests:
            result = instagram.verify_address(test['data'])
            self.assertEqual(test['valid'], result, test['name'])


if __name__ == '__main__':
    unittest.main()
