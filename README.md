# Themis
> Themis, (Greek: “Order”) in Greek religion, personification of justice, goddess of wisdom and good counsel, and the interpreter of the gods' will.

Themis is a collection of REST APIs, [Band Protocol data sources](https://docs.bandchain.org/whitepaper/terminology.html#data-sources) and [Band Protocol oracle scripts](https://docs.bandchain.org/whitepaper/terminology.html#oracle-scripts) that work together in order to make it possible for [Desmos](https://desmos.network) users to connect their profile to common social network profiles owned by them.

The connection process is defined as follows: 
1. The user signs their social network username with the private key associated with their Desmos profile. 
2. The signature is posted online, so that it can be reached by a public link ([example here](https://pastebin.com/raw/xz4S8WrW)).
3. The user posts the public link to the signature on one of the supported social networks listed below ([example Tweet here](https://twitter.com/ricmontagnin/status/1368883070590476292)). 
4. Data sources will use our APIs to get that data into the Band chain and verify their correctness (valid signature, and valid address). 
5. Once the data has been verified, the oracle scripts will then store it inside the Band chain as a proof of connection so that it can be later queried and verified. 

At the end of the process, we should have proved with sufficient certainty that: 
- the user possesses the private key associated to a Desmos profile (as they have been able to sign a message with it)
- the same user also possesses access to the social network accounts (as they have been able to post the link somewhere related to them)

So, we can conclude that the Desmos profile and the social network profile should be connected together. 

## Currently supported social networks
Currently, we support the following ways of verifying a profile with a social network.

### Twitter
- Using a tweet (if the profile is public)
- Using the profile biography (description), even if the profile itself is private 

## Components
### APIs
Our APIs are a wrapper around the different social networks APIs that we will use to get the data that has been posted online (eg. Twitter APIs to get a tweet or a profile's biography). Their code is present inside the `apis` folder.

To run these APIs, all you need to do is create a configuration file using the TOML format. Then, you can run the `main.go` function using the following command: 

```shell
go run main.go <path/to/config/file.toml>
```

#### Configuration
The configuration file should contain the following data: 

```toml
[twitter]
bearer="<Bearer token used to access the Twitter APIs>"
cache_file="<Absolute path to the cache file where requests will be cached>"
```

### Data sources
The data sources we use are simple Python scripts that perform the following operations:

1. They get the data from the social network using [our REST APIs](#apis).
2. They verify the validity of such data by making sure it contains all the proper data (signature, signed message, public key, and address) and that those data is valid .

Once they have performed all these operations, they return the first URL that contained a valid signature that they found. This is done to allow the oracle scripts (that will use these data sources) to then save the URL inside the chain. 

The code of different data sources is found inside the `data-sources` folder.

### Oracle scripts
The last part of our system are oracle scripts. These are responsible for calling data sources that will fetch the data from centralized social networks and check to make sure that data is valid. 

If the retrieved data is valid, then data sources will return a URL that is possible to use to get the verified data. Oracle scripts will then just need to make sure that enough data sources have returned a URL and have not errored. Once that requirement is satified, they will just store the valid URL inside the Band chain.

All the oracle scripts code is present inside the `oracle-scripts` folder. Such code is written in Rust and must be compatible with [OWasm](https://docs.rs/owasm/0.1.10/owasm/).

