import unittest
import github


class GitHubTest(unittest.TestCase):

    def test_get_data_from_gist(self):
        # Valid signature
        data = github.get_data_from_gist(github.CallData('RiccardoM', '720e0072390a901bb80e59fd60d7fded'))
        self.assertIsNotNone(data)

        # Invalid signature
        data = github.get_data_from_gist(github.CallData('RiccardoM', 'invalid_gist_id'))
        self.assertIsNone(data)


if __name__ == '__main__':
    unittest.main()
