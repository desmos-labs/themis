import unittest
import domain
import httpretty


class DomainTest(unittest.TestCase):

    def test_get_data_from_txt_record(self):
        jsons = [
            {
                'name': 'Valid JSON',
                'value': '{\"address\": \"470289c128641e77756e5a5bfaf6832d0e5c7211\",\"pub_key\": \"02fc6a0f6001262c38dc0d1ec34b476ced1c394db2927860fb0359f5ba4a9cd964\",\"signature\": \"e22a37f2a4e5c319bb460daf2e8113a4ab94f55687cbaa930364e3f20333e8810abbb08ee9d9236e74117088b7c214c8d62d524a55374e596e131073ecd5c113\",\"value\": \"676f66696e645f6d65\"}',
                'valid': True,
            },
            {
                'name': 'Valid URL',
                'value': 'https://api.go-find.me/proof/4873be0a74ac3f6ffc05d5e2c82614c8',
                'valid': True,
            },
            {
                'name': 'Invalid string',
                'value': 'google-site-verification=TGkpVO2zOmAs4jva1JtL1Jg3r6xHLwK8pc7x-I2_AUQ',
                'valid': False,
            },
        ]

        for json in jsons:
            result = domain.get_data_from_txt_record(json['value'])
            if json['valid'] is True:
                self.assertIsNotNone(result)
            else:
                self.assertIsNone(result)

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
            result = domain.validate_json(json['json'])
            self.assertEqual(json['valid'], result, json['name'])

    @httpretty.activate(verbose=True, allow_net_connect=False)
    def test_get_user_data(self):
        # Register fake HTTP call
        httpretty.register_uri(
            httpretty.GET,
            "https://themis.morpheus.desmos.network/nslookup/forbole.com",
            status=200,
            body='{"txt":[{"text":"{\\"address\\": \\"470289c128641e77756e5a5bfaf6832d0e5c7211\\",\\"pub_key\\": \\"02fc6a0f6001262c38dc0d1ec34b476ced1c394db2927860fb0359f5ba4a9cd964\\",\\"signature\\": \\"e22a37f2a4e5c319bb460daf2e8113a4ab94f55687cbaa930364e3f20333e8810abbb08ee9d9236e74117088b7c214c8d62d524a55374e596e131073ecd5c113\\",\\"value\\": \\"676f66696e645f6d65\\"}"},{"text":"google-site-verification=TGkpVO2zOmAs4jva1JtL1Jg3r6xHLwK8pc7x-I2_AUQ"}]}',
        )

        # Valid signature
        data = domain.get_user_data(domain.CallData('forbole.com'))
        self.assertIsNotNone(data)

        # Invalid signature
        httpretty.register_uri(
            httpretty.GET,
            "https://themis.morpheus.desmos.network/nslookup/forbole.com",
            status=404,
        )
        data = domain.get_user_data(domain.CallData('forbole.com'))
        self.assertIsNone(data)

    def test_verify_signature(self):
        tests = [
            {
                'name': 'Valid data',
                'valid': True,
                'data': domain.VerificationData(
                    '',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    '7269636d6f6e7461676e696e',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                ),
            },
            {
                'name': 'Invalid value',
                'valid': False,
                'data': domain.VerificationData(
                    '',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    'ricmontagni',
                ),
            },
            {
                'name': 'Invalid signature',
                'valid': False,
                'data': domain.VerificationData(
                    '',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'a00a7d5bd45e42615645fcaeb4d800af2704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    '7269636d6f6e7461676e696e',
                ),
            },
            {
                'name': 'Invalid pub key',
                'valid': False,
                'data': domain.VerificationData(
                    '',
                    '033024e9e0ad4f9305ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173',
                    '7269636d6f6e7461676e696e',
                ),
            },
        ]

        for test in tests:
            result = domain.verify_signature(test['data'])
            self.assertEqual(test['valid'], result, test['name'])

    def test_verify_address(self):
        tests = [
            {
                'name': 'Valid address',
                'valid': True,
                'data': domain.VerificationData(
                    '8902A4822B87C1ADED60AE947044E614BD4CAEE2',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d',
                    '7269636d6f6e7461676e696e',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173'
                ),
            },
            {
                'name': 'Invalid address',
                'valid': False,
                'data': domain.VerificationData(
                    '8902A4822B87C1ADED60AE947044E614BD4CAEE2',
                    '033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b41',
                    '7269636d6f6e7461676e696e',
                    'a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173'
                ),
            },
        ]

        for test in tests:
            result = domain.verify_address(test['data'])
            self.assertEqual(test['valid'], result, test['name'])


if __name__ == '__main__':
    unittest.main()
