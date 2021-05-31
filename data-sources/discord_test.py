import unittest
import discord


class DiscordTest(unittest.TestCase):

    def test_get_data_for_user(self):
        # Valid signature
        discord.ENDPOINT = "http://localhost:5000/discord"
        data = discord.get_user_data(discord.CallData('Riccardo Montagnin#5414'))
        self.assertIsNotNone(data)

        # Invalid signature
        data = discord.get_user_data(discord.CallData('RiccardoM'))
        self.assertIsNone(data)


if __name__ == '__main__':
    unittest.main()
