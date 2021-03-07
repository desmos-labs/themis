# Themis
> Themis, (Greek: “Order”) in Greek religion, personification of justice, goddess of wisdom and good counsel, and the interpreter of the gods' will.

Themis is a wrapper around some common social networks APIs, that has been created in order to be used by [Band Protocol data sources](https://docs.bandchain.org/whitepaper/terminology.html) during the process of verifying the ownership of social network accounts by a specific Desmos profile.  

The thought process is the following: 
1. The user signs their username with the private key associated with the Desmos profile. 
2. The signature is posted online, so that it can be reached by a public link.
3. The user posts the public link to the signature on one of the above mentioned sources (tweet, profile bio, etc). 
4. Data sources will use Themis to get that data into the Band chain. 
5. The data will be used by [Oracle scripts](https://docs.bandchain.org/whitepaper/terminology.html) to verify the authenticity and stored the validation result inside the Band chain itself. 

At the end of the process, we should have proved with sufficient certainty that: 
- the user possesses the private key associated to a Desmos profile 
- the same user also possesses access to the social network accounts (as they have been able to post the link somewhere related to them)

So, we can conclude that the Desmos profile and the social network profile should be connected together. 

## Currently supported social networks
Currently, we support the following ways of verifying a profile with a social network: 
- Twitter
   - Using a tweet (if the profile is public)
   - Using the profile biography (description), even if it is private 

## Run 
To run these APIs, all you need to do is create a configuration file using the TOML format. Then, you can run the `main.go` function using the following command: 

```shell
go run main.go <path/to/config/file.toml>
```

## Configuration
The configuration file should contain the following data: 

```toml
[twitter]
bearer="<Bearer token used to access the Twitter APIs>"
cache_file="<Absolute path to the cache file where requests will be cached>"
```