# Data sources
The data sources we use are simple Python scripts that perform the following operations:

1. They get the data from the social network using [our REST APIs](../apis/README.md).
2. They verify the validity of such data by making sure it contains all the proper fields (signature, signed message, public key, and address) and that data is valid.

Once they have performed all these operations, they return the first URL that contained a valid signature that they found. This is done to allow the oracle scripts (that will use these data sources) to then save the URL inside the chain.

## Uploading to Band Protocol
In order to work, each data source must be uploaded inside the Band Protocol blockchain. To do this, you can download the `bandcli` executable and then run the following command: 

```shell
$ bandcli tx oracle create-data-source
# Create a new data source that will be used by oracle scripts.
# Usage:
#   bandcli tx oracle create-data-source (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner]) [flags]
```

Example: 

```shell
$ bandcli tx oracle create-data-source \
  --name themis-twitter \
  --description "Data source allowing to verify a Twitter account" \
  --script ./twitter.py \
  --owner <your_address> 
```

#### Note  
Please make sure you **always** specify an owner of the data source using the `--owner` flag. This will make it possible for you to edit the data source in the future if you want so. Not specifying an owner will result in an immutable data source. 