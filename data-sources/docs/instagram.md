# Instagram

Users can connect their Instagram account by posting a link to the verification data inside the biography description.  
The verification data must be a JSON object formed as described inside ["Verification data"](../../README.md#verification-data).

## Example biography

An example biography can be found [here](https://www.instagram.com/test_desmos_user).

## Verification process

The verification process on Instagram is made of the following steps:

1. The user uploads their verification data to an online storage (eg. PasteBin).
2. The user links the online storage URL to their account by putting it inside their account biography.
3. The user requests Themis to cache their user profile data by an access token with the permission requesting for their profile.
4. The user performs a Desmos transaction telling that they want to link the Instagram account to the Desmos one, and provides the proper call data.

Once that's done, what will happen is the following:

1. Desmos will send the call data to Band Protocol, asking to get the Instagram user username and the signature provided by the user.
2. Band Protocol will call the appropriate data source that will use our APIs to get the data from the account biography.
3. Once downloaded, the data source will check the validity of the data and return to Desmos the user username and signature.
4. If the Instagram user username matches the one provided by the user, and the signature is valid against the user public key, the Desmos and Instagram account will be linked together.  

## Data source call data

When asking to verify the ownership of a Instagram account, the data source call data must be a JSON object formed as follows:

```json
{
  "username": "<Instagram user username to be verified>"
}
```

Example:

```json
{
  "username":"test_desmos_user"
}
```

Hex encoded:
```
7b22757365726e616d65223a22746573745f6465736d6f735f75736572227d
```

Example execution:

```shell
python Instagram.py 7b22757365726e616d65223a22746573745f6465736d6f735f75736572227d
```
