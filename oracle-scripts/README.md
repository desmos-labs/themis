# Oracle scripts
Oracle scripts are responsible for calling [data sources](../data-sources) that will fetch the data from centralized social networks and check to make sure that data is valid.

If the retrieved data is valid, then the data sources will return the centralized social network username and the hex-encoded signature. Oracle scripts will then just need to make sure that enough data sources have returned without error. Once that requirement is satisfied, they will just store the data inside the Band chain.

Since oracle scripts are executed on-chain, the code is written in Rust and must be compatible with [OWasm](https://docs.rs/owasm/0.1.10/owasm/).

## Customization
Before uploading the oracle script inside the Band Protocol blockchain, you need to customize it specifying the ID of the [data sources](../data-sources) that you want to be called by the script. 

To do this, you can edit the constants inside the `scr/script.rs` file: 

```rust
const DATA_SOURCE_TWITTER: i64 = 49;
const DATA_SOURCE_GITHUB: i64 = 68;
const DATA_SOURCE_DISCORD: i64 = 80;
// Other sources
```

The IDs should be the ones of the data sources you want to call. To get them, you can simply use the [Band Protocol explorer](https://cosmoscan.io/data-sources). 

#### Note
If the ID displayed inside the block explorer is `#D15`, you only need to take the `15`, excluding the `#D` part.

## Uploading to Band Protocol
The first thing you have to do when you want to upload your oracle script to the Band Protocol chain is to compile it. To do this, you can use the following command: 

```shell
RUSTFLAGS='-C link-arg=-s' cargo build --target wasm32-unknown-unknown --release
```

Then, the compiled code must be uploaded inside the Band Protocol blockchain. To do this, you can download the `bandd` executable and then run the following command:

```shell
$ bandd tx oracle create-oracle-script --help
# Create a new oracle script that will be used by data requests.
# Usage:
#  bandd tx oracle create-oracle-script (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner]) (--schema [schema]) (--url [source-code-url]) [flags
```

Example:

```shell
$ bandd tx oracle create-oracle-script \
  --name desmos-themis \
  --description "Oracle script allowing to verify a Desmos account" \
  --script target/wasm32-unknown-unknown/release/themis_oracle_script.wasm \
  --treasury <your_address> \
  --owner <your_address> 
```

#### Note
Please make sure you **always** specify an owner of the oracle script using the `--owner` flag. This will make it possible for you to edit the script in the future if you want so. Not specifying an owner will result in an immutable owner.

## Editing the oracle script
If you want to edit an oracle script, you should use the `edit-oracle-script` command: 

```shell
$ bandd tx oracle edit-oracle-script
# Edit an existing oracle script that will be used by data requests.
# Usage:
#  bandd tx oracle edit-oracle-script [id] (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner]) (--schema [schema]) (--url [source-code-url]) [flags]

```

Example: 

```shell
$ bandd tx oracle edit-oracle-script 32 \
  --url https://raw.githubusercontent.com/desmos-labs/themis/main/oracle-scripts/src/script.rs \
  --owner $(bandd keys show jack -a) \
  --from jack
```

## Calling the oracle script
In order to properly call an oracle script, you have to run the following command: 

```shell
$ bandcli tx oracle request
# Make a new request via an existing oracle script with the configuration flags.
# Usage:
#  bandd tx oracle request [oracle-script-id] [ask-count] [min-count] (-c [calldata]) (-m [client-id]) (--prepare-gas=[prepare-gas] (--execute-gas=[execute-gas])) (--fee-limit=[fee-limit]) [flags]
```

Example:
```shell
$ bandd tx oracle request 32 7 4 \
  -c 0000000574776565740000001331333932303333353835363735333137323532 \
  --gas 600000 \
  --from jack
```

The call data must be OBI encoded. To easily get it, you can use the `test_obi_encode` test inside the `script.rs` file. Simply change the attributes of the `input` variable and then run the test. The output will be the OBI encoded call data that you must use to perform your oracle request. 

## Getting the result from a script
In order to get the result of the execution of an oracle script, you can use the following command: 

```shell
bandd q oracle request
# Usage:
#  bandcli query oracle request [id] [flags]
```

To get the request id, you can simply query the transaction that you used to make such request. Inside the output you will find the `request` type event and then get the `id` attribute:

```
height: 7862621
txhash: 0AD75B809F185B9B101222150072A94A701BA8DD6754C1E64F8E84071FE2D194
logs:
- msgindex: 0
  log: ""
  events:
  - type: message
    attributes:
    - key: action
      value: request
  - type: request
    attributes:
    - key: id
      value: "3585983" <--- Request id here
```

From here, you can then call the command with that request id:

```shell
$ bandcli q oracle request 3585983
```

Now, you need to get the `result.result` value: 

```
reports: []
request: null
result:
  ans_count: "10"
  ask_count: "10"
  calldata: AAAABmdpdGh1YgAAAIo3QjIyNzU3MzY1NzI2RTYxNkQ2NTIyM0EyMjUyNjk2MzYzNjE3MjY0NkY0RDIyMkMyMjY3Njk3Mzc0NUY2OTY0MjIzQTIyMzczMjMwNjUzMDMwMzczMjMzMzkzMDYxMzkzMDMxNjI2MjM4MzA2NTM1Mzk2NjY0MzYzMDY0Mzc2NjY0NjU2NDIyN0Q=
  client_id: desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu-github-RiccardoM
  min_count: "6"
  oracle_script_id: "32"
  request_id: "148493"
  request_time: "1622454766"
  resolve_status: RESOLVE_STATUS_SUCCESS
  resolve_time: "1622454772"
  result: AAAAgDkyYjFiZmFjOTQ0YmUwOTZlOGI4MDA4MTNhMGQ5OGFjODU1ODlkNDI3MjA2YTJiMTlkMDAxMzE3MmRkNGViMDc0YjJiY2I2NTViNTJkMmJiMGZmNjkxNDVhZjI0YWM5OWRjOWQzY2Y1ZTk2OWVmMTg5NzQ5YjM1ZjM5MDU5NWM2AAAACVJpY2NhcmRvTQ==
```

To decode the output, you can use the `test_obi_decode` inside the `script.rs` file, by replacing the integer vector with your output. 
