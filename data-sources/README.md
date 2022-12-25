# Data sources
The data sources we use are simple Python scripts that perform the following operations:

1. They get the data from the social network.
2. They verify the validity of such data by making sure it contains all the proper fields (signature, signed message, public key, and address) and that data is valid.

Once they have performed all these operations, they return the social network username as well as the signature value. This is done to allow the oracle scripts (that will use these data sources) to then save such data inside the chain.

## Call data
Each data source must be called passing to it the so-called `call data`. This is a hex-encoded JSON object containing all the arguments values that such data source will use to run the verifications. 

We use this method of passing the arguments because different data sources might want different data to run properly. Using a hex encoded JSON object we can simply create such object anywhere and then send it to the oracle script that will just forward it to the data source. 

If you want to know more about what kind of different data each source requires, please read inside the [_supported apps_ folder](docs).

## Uploading to Band Protocol
In order to work, each data source must be uploaded inside the Band Protocol blockchain. To do this, you can download the `bandd` executable and then run the following command: 

```shell
$ bandd tx oracle create-data-source
# Create a new data source that will be used by oracle scripts.
# Usage:
#  bandd tx oracle create-data-source (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner]) (--treasury [treasury]) (--fee [fee]) [flags]
```

Example: 

```shell
$ bandd tx oracle create-data-source \
  --name themis-twitter \
  --description "Data source allowing to verify a Twitter account" \
  --script ./twitter.py \
  --treasury <your_address> \
  --owner <your_address> 
```

#### Note  
Please make sure you **always** specify an owner of the data source using the `--owner` flag. This will make it possible for you to edit the data source in the future if you want so. Not specifying an owner will result in an immutable data source. 
