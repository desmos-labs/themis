#!/usr/bin/env python3
import json
import sys
import requests
import re
from typing import Optional
import cryptography.hazmat.primitives.asymmetric.utils as crypto
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import hashes
from ripemd.ripemd160 import ripemd160
import hashlib

METHOD_TWEET = "tweet"
METHOD_PROFILE = "bio"
TYPES = [METHOD_TWEET, METHOD_PROFILE]

ENDPOINT = "https://themis.morpheus.desmos.network/twitter"
HEADERS = {"Content-Type": "application/json"}


class CallData:
    """
    Contains the data that has been used to call the script
    """

    def __init__(self, method: str, value: str):
        self.method = method
        self.value = value


class VerificationData:
    """
    Contains the data needed to verify the proof submitted by the user.
    """

    def __init__(self, address: str, pub_key: str, value: str, signature: str):
        self.address = address
        self.pub_key = pub_key
        self.signature = signature
        self.value = value


def get_data_from_tweet(tweet: str):
    """
    Returns the username of the creator and all the URLs found inside the tweet having the given id.
    :param tweet: Id of the Tweet to be fetched
    :return: Username of the Tweet's creator and all the found URLs
    """
    result = requests.request("GET", f"{ENDPOINT}/tweets/{tweet}", headers=HEADERS).json()
    return result['author']['username'], re.findall(r'(https?://[^\s]+)', result['text'])


def get_data_from_bio(user: str):
    """
    Returns the username and all the URLs that are found inside the bio of the user having the given username.
    :param user: Username of the user for whom to check the bio.
    :return: List of URLs found inside the bio of the user.
    """
    result = requests.request("GET", f"{ENDPOINT}/users/{user}", headers=HEADERS).json()
    return result['username'], re.findall(r'(https?://[^\s]+)', result['bio'])


def get_signature_from_url(url: str) -> Optional[VerificationData]:
    """
    Tries getting the signature object linked to the given URL.
    :param url: URL that should contain the signature object.
    :return: A dictionary containing 'valid' to tell whether the search was valid, and an optional 'data' containing
    the signature object.
    """
    try:
        result = requests.request("GET", url, headers=HEADERS).json()
        if validate_json(result):
            return VerificationData(
                result['address'],
                result['pub_key'],
                result['value'],
                result['signature'],
            )
        else:
            return None
    except ValueError:
        return None


def validate_json(json: dict) -> bool:
    """
    Tells whether or not the given JSON is a valid signature JSON object.
    :param json: JSON object to be checked.
    :return: True if the provided JSON has a valid signature schema, or False otherwise.
    """
    return all(key in json for key in ['value', 'pub_key', 'signature', 'address'])


def verify_signature(data: VerificationData) -> bool:
    """
    Verifies the signature using the given pubkey and value.
    :param data: Data used to verify the signature.
    :return True if the signature is valid, False otherwise
    """
    if len(data.signature) != 128:
        return False

    try:
        # Create signature for dss signature
        (r, s) = int(data.signature[:64], 16), int(data.signature[64:], 16)
        sig = crypto.encode_dss_signature(r, s)

        # Create public key instance
        public_key = ec.EllipticCurvePublicKey.from_encoded_point(ec.SECP256K1(), bytes.fromhex(data.pub_key))

        # Verify the signature
        public_key.verify(sig, bytes.fromhex(data.value), ec.ECDSA(hashes.SHA256()))
        return True
    except Exception:
        return False


def verify_address(data: VerificationData) -> bool:
    """
    Verifies that the given address is the one associated with the provided HEX encoded compact public key.
    :param data: Data used to verify the address
    """
    s = hashlib.new("sha256", bytes.fromhex(data.pub_key)).digest()
    r = ripemd160(s).hex()
    return data.address.upper() == r.upper()


def check_values(values: dict) -> CallData:
    """
    Checks the validity of the given dictionary making sure it contains the proper data.
    :param values: Dictionary that should be checked.
    :return: A CallData instance.
    """
    if "method" not in values:
        raise Exception("Missing 'method' value")

    if "value" not in values:
        raise Exception("Missing 'value' value")

    return CallData(values["method"], values["value"])


def main(args: str):
    """
    Gets the signature data from Twitter, from either a Tweet or a profile biography.

    :param args Hex encoded JSON object containing the arguments to be used during the execution.
    In order to be valid, the encoded JSON object must contain two fields:

    1. "method", which represents the verification method to be used;
    2. "value", which represents the value associated with the verification method.

    The two possible values for "method" are:
    - "tweet" if the link is provided inside a public tweet.
       In this case, "value" must represent a valid Tweet ID.

    - "bio" if the link is provided inside the user's profile biography.
       In this case, "value" must be a valid Twitter username.

    Failing to provide any of these value will result in the wrong output being returned.

    In both cases, the link must be public and its content must be a JSON object formed as follows:

    ```json
    {
      "address": "Hex encoded address of the signer",
      "pub_key": "Hex encoded public key that has been used to sign the value",
      "value": "Hex encoded value that has been signed",
      "signature": "Hex encoded Secp256k1 signature"
    }
    ```

    Example argument value:
    7B226D6574686F64223A227477656574222C2276616C7565223A2231333932303333353835363735333137323532227D

    This is the hex encoded representation of the following JSON object:

    ```json
    {
      "method":"tweet",
      "value":"1392033585675317252"
    }
    ```

    :param args: JSON encoded parameters used during the execution.
    :return The signed value and the signature as a single comma separated string.
    :raise Exception if anything is wrong during the process. This can happen if:
            1. The tweet or profile bio does not contain any valid URL that link to a valid signature data object
            2. The provided signature is not valid
            3. The provided address is not linked to the provided public key
    """

    decoded = bytes.fromhex(args)
    json_obj = json.loads(decoded)
    call_data = check_values(json_obj)

    verification_method = call_data.method
    if verification_method not in TYPES:
        raise Exception(f"Invalid verification method: {verification_method}")

    # Get the URLs to check inside the tweet or the bio
    username = ""
    urls = []
    if verification_method == METHOD_TWEET:
        username, urls = get_data_from_tweet(call_data.value)
    elif verification_method == METHOD_PROFILE:
        username, urls = get_data_from_bio(call_data.value)

    if len(urls) == 0:
        raise Exception(f"No URL found inside {verification_method}")

    # Find the signature following the URLs
    data = None
    for url in urls:
        result = get_signature_from_url(url)
        if result is not None:
            data = result
            break

    if data is None:
        raise Exception(f"No valid signature data found inside {verification_method}")

    # Verify the signature
    signature_valid = verify_signature(data)
    if not signature_valid:
        raise Exception("Invalid signature")

    # Verify the address
    address_valid = verify_address(data)
    if not address_valid:
        raise Exception("Invalid address")

    return f"{data.value},{data.signature},{username}"


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
