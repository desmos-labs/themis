# Oracle scripts
Oracle scripts are responsible for calling [data sources](../data-sources/README.md) that will fetch the data from centralized social networks and check to make sure that data is valid.

If the retrieved data is valid, then data sources will return a URL that is possible to use to get the verified data. Oracle scripts will then just need to make sure that enough data sources have returned a URL and have not errored. Once that requirement is satisfied, they will just store the valid URL inside the Band chain.

Since oracle scripts are executed on-chain, the code is written in Rust and must be compatible with [OWasm](https://docs.rs/owasm/0.1.10/owasm/).

## Customization
Before uploading the oracle script inside the Band Protocol blockchain, you need to customize it specifying the ID of the [data source](../data-sources/README.md) that you want to be called by the script. 

To do this, you can edit the following constant inside the `scr/script.rs` file: 

```rust
const DESMOS_THEMIS_DS: i64 = 49;
```

This ID should be the one of the data source you want to call. To get it, you can simply use the [Band Protocol explorer](https://cosmoscan.io/data-sources) and get it from there. 

#### Note
If the ID displayed inside the block explorer is `#D15`, you only need to take the `15`, excluding the `#D` part.

## Uploading to Band Protocol
In order to be executable, each oracle script must be uploaded inside the Band Protocol blockchain. To do this, you can download the `bandcli` executable and then run the following command:

```shell
$ bandcli tx oracle create-oracle-script --help
# Create a new oracle script that will be used by data requests.
# Usage:
#   bandcli tx oracle create-oracle-script (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner]) (--schema [schema]) (--url [source-code-url]) [flags]
```

Example:

```shell
$ bandcli tx oracle create-oracle-script \
  --name themis-twitter \
  --description "Oracle script allowing to verify a Twitter account" \
  --script ./target/wasm32-unknown-unknown/themis_oracle_script.wasm \
  --owner <your_address> 
```

#### Note
Please make sure you **always** specify an owner of the oracle script using the `--owner` flag. This will make it possible for you to edit the script in the future if you want so. Not specifying an owner will result in an immutable owner. 