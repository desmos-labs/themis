import unittest
import urllib.parse

import discord
import httpretty


class DiscordTest(unittest.TestCase):

    @httpretty.activate(verbose=True, allow_net_connect=False)
    def test_get_data_for_user(self):
        # Register fake HTTP call
        httpretty.register_uri(
            httpretty.GET,
            "https://themis.morpheus.desmos.network/discord/Riccardo%20Montagnin%235414",
            status=200,
            body='{"address":"8902A4822B87C1ADED60AE947044E614BD4CAEE2","pub_key":"033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d","value":"Riccardo Montagnin#5414","signature":"d10db146bb4d234c5c1d2bc088e045f4f05837c690bce4101e2c0f0c6c96e1232d8516884b0a694ee85e9c9da51be74966886cbb12af4ad87e5336da76d75cfb"}',
        )

        # Valid signature
        data = discord.get_user_data(discord.CallData('Riccardo Montagnin#5414'))
        self.assertIsNotNone(data)

        # Invalid signature
        httpretty.register_uri(
            httpretty.GET,
            "https://themis.morpheus.desmos.network/discord/RiccardoM",
            status=404,
        )
        data = discord.get_user_data(discord.CallData('RiccardoM'))
        self.assertIsNone(data)


if __name__ == '__main__':
    unittest.main()
