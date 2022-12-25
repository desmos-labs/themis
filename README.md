# Themis
> Themis, (Greek: “Order”) in Greek religion, personification of justice, goddess of wisdom and good counsel, and the interpreter of the gods' will.

Themis is a collection of REST APIs, [Band Protocol data sources](https://docs.bandchain.org/whitepaper/terminology.html#data-sources) and [Band Protocol oracle scripts](https://docs.bandchain.org/whitepaper/terminology.html#oracle-scripts) that work together in order to make it possible for [Desmos](https://desmos.network) users to connect their profile to common social network profiles owned by them.

The connection process is defined as follows: 
1. The user signs their social network username with the private key associated with their Desmos profile. 
2. The signature is posted online from the social network account, so that it can be reached publicly (eg. inside a Tweet, a Gist, etc.).
3. Data sources will use our APIs to get that data into the Band chain and verify their correctness (valid signature, and valid address). 
4. Once the data has been verified, the oracle scripts will then store it inside the Band chain as a proof of connection so that it can be later queried and verified. 

At the end of the process, we should have proved with sufficient certainty that: 
1. the user possesses the private key associated to a Desmos profile (as they have been able to sign a message with it); and
2. the same user also possesses access to the social network accounts (as they have been able to post the link using it).

So, we can conclude that the Desmos profile and the social network profile should be connected together.

## Process flow
Following you can find the graphical representation of the overall process flow: 

![](img/flow.png)

If you want to know more about the individual components, please read the `README.md` files inside each folder.

### Verification data
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

#### Validity checks
Once the data has been retrieved from the proper location (eg. tweet, Gist, etc.), the following checks will be performed:

1. The `signature` is a valid signature of the `value` made using the private key associated with the `pub_key`.
2. The `address` is a valid address associated with the provided `pub_key`.


## Documentation 
If you want to know which application we support and how everything works, please refer to the [_.docs_ folder](docs).
