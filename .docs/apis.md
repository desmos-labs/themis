# APIs
Our APIs are a wrapper around the different social networks APIs that we will use to get the data that has been posted online (eg. Twitter APIs to get a tweet or a profile's biography).

To run these APIs, all you need to do is create a configuration file using the TOML format. Then, you can run the `main.go` function using the following command:

```shell
go run main.go <path/to/config/file.toml>
```

## Configuration
The configuration file should contain the following data:

```toml
[apis]
port = <Port on which to run the APIs>

[twitter]
bearer = "<Bearer token used to access the Twitter APIs>"
cache_file = "<Absolute path to the cache file where requests will be cached>"

[discord]
store_folder_path = "<Path to the folder where data will be stored>"
bot_pub_key_path = "<Path to the public key file contining Hephaestus' public key>"

[twitch]
client_id = "<Twitch APIs client id>"
client_secret = "<Twitch APIs client secret>"
```
