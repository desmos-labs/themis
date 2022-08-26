# Youtube
Users can connect their youtube account by posting a link to the verification data via:

- the account channel description 

In both cases, the verification data must be a JSON object formed as described inside ["Verification data"](../verification-data.md).

## Example channel
An example channel can be found [here](https://www.youtube.com/channel/UCvu-RHbVTuLvF3HZqzEV-ZA/about). 

## Verification process
The verification process on Youtube is made of the following steps:

1. The user uploads their verification data to an online storage (eg. PasteBin). 
2. The user links the online storage URL to their account by putting it inside their account channel description.
3. The user performs a Desmos transaction telling that they want to link the Youtube account to the Desmos one, and provides the proper call data. 

Once that's done, what will happen is the following: 

1. Desmos will send the call data to Band Protocol, asking to get the Youtube user id and the signature provided by the user. 
2. Band Protocol will call the appropriate data source that will use our APIs to get the data from the account channel description. 
3. Once downloaded, the data source will check the validity of the data and return to Desmos the user id and signature. 
4. If the Youtube user id matches the one provided by the user, and the signature is valid against the user public key, the Desmos and Youtube account will be linked together.  



## Data source call data
When asking to verify the ownership of a Youtube account, the data source call data must be a JSON object formed as follows: 

```json
{
  "user_id": "<Youtube user id to be verified",
}
```

Example: 
```json
{
  "user_id":"vu-RHbVTuLvF3HZqzEV-ZA"
}
```

Hex encoded:
```
7ba202022757365725f6964223a2276752d5248625654754c764633485a717a45562d5a4122a7d
```

Example execution: 

```shell
python youtube.py 7b22757365725f6964223a2276752d5248625654754c764633485a717a45562d5a41227d
```
