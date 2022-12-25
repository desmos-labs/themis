# GitHub
Users can connect their GitHub account by posting a **public Gist**. Such Gist must contain the JSON object as defined inside ["Verification data"](../../README.md#verification-data).

## Example Gist
An example Gist can be found [here](https://gist.github.com/RiccardoM/720e0072390a901bb80e59fd60d7fded).

## Verification process
The verification process of a GitHub account is made of two steps: 

1. The user creates a public Gist containing the verification data.
2. The user performs a Desmos transaction telling that they want to verify their GitHub account and providing the proper call data.

Once that's done, the following will happen: 

1. Desmos will ask Band to get the data posted inside the Gist. 
2. Band Protocol will use the appropriate data source to get the raw version of the Gist. 
3. Once retrieved, the data is going to be checked for validity and, if valid, the GitHub username and signature are returned. 
4. Desmos will get the returned signature and username, and perform the final checks on them. If everything is valid, it will create the GitHub link.

![](../../img/github.png)

### About Gist validity  
Please note that since the Gist is fetched using the HTTPS APIs, only the **first version** of the Gist will always be used. This means that if you fail in providing the correct value, you should create a **new Gist** instead of editing an existing one. Editing an old one and asking for the verification will fail.

## Data source call data
When asking to verify a GitHub username, the call data must be a JSON composed as follows: 

```json
{
  "username": "<GitHub username to be verified>",
  "gist_id": "<Gist ID used for the verification>"
}
```

Example: 
```json
{
  "username":"RiccardoM",
  "gist_id":"720e0072390a901bb80e59fd60d7fded"
}
```

Hex encoded: 
```
7B22757365726E616D65223A22526963636172646F4D222C22676973745F6964223A223732306530303732333930613930316262383065353966643630643766646564227D
```

Example execution:
```shell
python github.py 7B22757365726E616D65223A22526963636172646F4D222C22676973745F6964223A223732306530303732333930613930316262383065353966643630643766646564227D
```
