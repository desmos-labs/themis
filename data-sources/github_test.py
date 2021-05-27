import unittest
import github


class GitHubTest(unittest.TestCase):

    def test_get_data_from_gist(self):
        # Valid signature
        data = github.get_data_from_gist(github.CallData('RiccardoM', 'f227ec0ddd5bf931b0b00b76b93ff149'))
        self.assertIsNotNone(data)

        # Invalid signature
        data = github.get_data_from_gist(github.CallData('RiccardoM', 'invalid_gist_id'))
        self.assertIsNone(data)


if __name__ == '__main__':
    unittest.main()
