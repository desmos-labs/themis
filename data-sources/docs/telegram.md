# Telegram
Users can connect their Telegram account via [_Hephaestus_, our Telegram bot](https://github.com/desmos-labs/hephaestus). To do so all they have to do is use the `/connect` command and provide the following data: 

- the hex encoded address or their Desmos account;
- the hex encoded public key associated with their Desmos account;
- the hex encoded signature of their Telegram username. 

Once that is provided, the bot will take care of uploading the data to [out APIs](../../apis) so that it can later be fetched from the appropriate data source.

## Example command
Following an example command: 

```
/connect testnet {"address":"71b0310267b49279116835ed35791c24c110012f","pub_key":"0203233fabd69a1b7a90bb968a0ab66e3af61989f65cf0bc1f8e9518740a302f1f","signature":"c12605456b8652df655bb43d0166586dfc0c5d758b03f127ca6b027d0ec140ca29b9569a20c9b78b72e13d15c1a7fa0b142dc0e624f3f51ef76bd94e55345d2a","value":"746573745f757365725f6465736d6f73"}
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

