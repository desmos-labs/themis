# Telegram
Users can connect their Telegram account via [_Hephaestus_, our Telegram bot](https://github.com/desmos-labs/hephaestus). To do so all they have to do is use the `/connect` command and provide the following data: 

- the hex encoded address or their Desmos account;
- the hex encoded public key associated with their Desmos account;
- the hex encoded signature of their Telegram username. 

Once that is provided, the bot will take care of uploading the data to [out APIs](../apis.md) so that it can later be fetched from the appropriate data source.

## Example command
Following an example command: 

```
/connect 8902A4822B87C1ADED60AE947044E614BD4CAEE2 033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d d10db146bb4d234c5c1d2bc088e045f4f05837c690bce4101e2c0f0c6c96e1232d8516884b0a694ee85e9c9da51be74966886cbb12af4ad87e5336da76d75cfb
```

## Verification process
The verification process is more complicated than other social networks due to the fact that Telegram does not allow anyone to post anything publicly. For this reason, the only way to upload some public data is using a Telegram bot. In our case we are using [_Hephaestus_, our official bot](https://github.com/desmos-labs/hephaestus). 

That being said, the verification process is the following: 

1. The user uses the `/connect` command to send the bot the required data (address, public key and signature). 
2. Hephaestus signs the data and uploads them to the Themis APIs.  
   The signature is created in order to avoid anyone maliciously sending false data to the APIs.
3. Once that's done, the user will receive from the bot the command that they need to run on Desmos. 
4. The user performs the command and starts the verification process.

Now, the following will happen:

1. Desmos will ask Band Protocol the data uploaded by the user. 
2. Band will use the data source to get the data from the Themis APIs. 
3. Once downloaded, the data source will verify the uploaded data and return the username as well as the signature. 
4. Desmos will verify the signature and the username, and create the link if everything is correct. 

## Data source call data
In order to call the data source, the call data must be a JSON object formed as follows: 

```json
{
  "username": "<Telegram username of the user>"
}
```

Example:
```json
{
  "username": "test_user_desmos"
}
```

Hex encoded:
```
207B22757365726E616D65223A22746573745F757365725F6465736D6F73227D
```

Example execution:
```shel
python telegram.py 207B22757365726E616D65223A22746573745F757365725F6465736D6F73227D
```

