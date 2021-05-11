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
The first thing you have to do when you want to upload your oracle script to the Band Protocol chain is to compile it. To do this, you can use the following command: 

```shell
RUSTFLAGS='-C link-arg=-s' cargo build --target wasm32-unknown-unknown --release
```

Then, the compiled code must be uploaded inside the Band Protocol blockchain. To do this, you can download the `bandcli` executable and then run the following command:

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
  --script target/wasm32-unknown-unknown/release/themis_oracle_script.wasm \
  --owner <your_address> 
```

#### Note
Please make sure you **always** specify an owner of the oracle script using the `--owner` flag. This will make it possible for you to edit the script in the future if you want so. Not specifying an owner will result in an immutable owner.

## Editing the oracle script
If you want to edit an oracle script, you should use the `edit-oracle-script` command: 

```shell
$ bandcli tx oracle edit-oracle-script
# Edit an existing oracle script that will be used by data requests.
# Usage:
#   bandcli tx oracle edit-oracle-script [id] (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner]) (--schema [schema]) (--url [source-code-url]) [flags]
```

Example: 

```shell
$ bandcli tx oracle edit-oracle-script 32 \
  --url https://raw.githubusercontent.com/desmos-labs/themis/main/oracle-scripts/src/script.rs \
  --owner $(bandcli keys show jack -a) \
  --from jack
```

## Calling the oracle script
In order to properly call an oracle script, you have to run the following command: 

```shell
$ bandcli tx oracle request
# Make a new request via an existing oracle script with the configuration flags.
# Usage:
#   bandcli tx oracle request [oracle-script-id] [ask-count] [min-count] (-c [calldata]) (-m [client-id]) [flags]
```

Example:
```shell
$ bandcli tx oracle request 32 7 4 \
  -c 0000000574776565740000001331333932303333353835363735333137323532 \
  --gas 600000 \
  --from jack
```

The call data must be OBI encoded. To easily get it, you can use the `test_obi_encode` test inside the `script.rs` file. Simply change the attributes of the `input` variable and then run the test. The output will be the OBI encoded call data that you must use to perform your oracle request. 

## Getting the result from a script
In order to get the result of the execution of an oracle script, you can use the following command: 

```shell
bandcli q oracle request
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

Now, you need to get the `responsepacketdata.result` value: 

```
request:
  oraclescriptid: 32
  mincount: 4
  requestheight: 7862621
  requesttime: 2021-05-11T08:37:20.754701045Z
  clientid: ""
    result:
      responsepacketdata:
        clientid: ""
        requestid: 3585983
        anscount: 7
        requesttime: 1620722240
        resolvetime: 1620722247
        resolvestatus: 1
        result:  <--- Get this
        - 1
        - 0
        - 0
        - 0
        - 24
        - 104
        - 116
        - 116
        - 112
        - 115
        - 58
        - 47
        - 47
        - 116
        - 46
        - 99
        - 111
        - 47
        - 117
        - 68
        - 50
        - 51
        - 72
        - 103
        - 83
        - 76
        - 74
        - 87
        - 10
```

To decode the output, you can use the `test_obi_decode` inside the `script.rs` file, by replacing the integer vector with your output. 