# Verification data
When proving the ownership of a centralized social network account, the verification data that should be posted using that account must be a valid JSON object formed as follows: 

```json
{
  "address":"<Hex encoded Desmos address>",
  "pub_key":"<Hex encoded Desmos public key>",
  "value":"<Plain text username>",
  "signature":"<Hex encoded signature of the username>"
}
```

Example: 
```json
{
  "address":"8902A4822B87C1ADED60AE947044E614BD4CAEE2",
  "pub_key":"033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d",
  "value":"RiccardoM",
  "signature":"92b1bfac944be096e8b800813a0d98ac85589d427206a2b19d0013172dd4eb074b2bcb655b52d2bb0ff69145af24ac99dc9d3cf5e969ef189749b35f390595c6"
}
```

## Validity checks
Once the data has been retrieved from the proper location (eg. tweet, Gist, etc.), the following checks will be performed: 

1. The `signature` is a valid signature of the `value` made using the private key associated with the `pub_key`. 
2. The `address` is a valid address associated with the provided `pub_key`.
