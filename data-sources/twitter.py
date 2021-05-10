#!/usr/bin/env python3

import sys
import requests
import re
from typing import Optional
import cryptography.hazmat.primitives.asymmetric.utils as crypto
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import hashes
import hashlib

TYPE_TWEET = "tweet"
TYPE_PROFILE = "profile"
TYPES = [TYPE_TWEET, TYPE_PROFILE]

ENDPOINT = "https://themis.morpheus.desmos.network/twitter"
HEADERS = {"Content-Type": "application/json"}


class VerificationData:
    """
    Contains the data needed to verify the proof submitted by the user.
    """

    def __init__(self, address: str, pub_key: str, value: str, signature: str):
        self.address = address
        self.pub_key = pub_key
        self.signature = signature
        self.value = value


def get_urls_from_tweet(tweet: str) -> [str]:
    """
    Returns all the URLs that are found inside the tweet having the given id.
    :param tweet: Id of the Tweet to be fetched
    :return: List of URLs that are found inside the tweet
    """
    url = f"{ENDPOINT}/tweets/{tweet}"
    result = requests.request("GET", url, headers=HEADERS).json()
    return re.findall(r'(https?://[^\s]+)', result['text'])


def get_urls_from_bio(user: str) -> [str]:
    """
    Returns all the URLs that are found inside the bio of the user having the given username.
    :param user: Username of the user for whom to check the bio.
    :return: List of URLs found inside the bio of the user.
    """
    url = f"{ENDPOINT}/users/{user}"
    result = requests.request("GET", url, headers=HEADERS).json()
    return re.findall(r'(https?://[^\s]+)', result['bio'])


def validate_json(json: dict) -> bool:
    """
    Tells whether or not the given JSON is a valid signature JSON object.
    :param json: JSON object to be checked.
    :return: True if the provided JSON has a valid signature schema, or False otherwise.
    """
    return all(key in json for key in ['value', 'pub_key', 'signature', 'address'])


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
        public_key.verify(sig, str.encode(data.value), ec.ECDSA(hashes.SHA256()))
        return True
    except Exception:
        return False


def verify_address(data: VerificationData) -> bool:
    """
    Verifies that the given address is the one associated with the provided HEX encoded compact public key.
    :param data: Data used to verify the address
    """
    s = hashlib.new("sha256", bytes.fromhex(data.pub_key)).digest()
    r = hashlib.new("ripemd160", s).digest()
    return data.address.upper() == r.hex().upper()


def main(type: str, value: str):
    """
    Gets the signature data from Twitter, either a tweet or a profile bio.
    The two possible values for "type" are:
    - "tweet" if the link is provided inside a public tweet
    - "profile" if the link is provided inside the user's profile biography

    The value must then be either the id of the tweet, or the username of the user.
    Failing to provide any of these value will result in the wrong output being returned.

    To be valid, a signature should be linked using a publicly available link, and should be composed as follows:

    {
      "address": "Bech 32 address os the signer",
      "pub_key": "Hex encoded public key that has been used to sign the value",
      "value": "Value that has been signed",
      "signature": "Hex encoded Secp256k1 signature"
    }

    :param type: Type of the verification that should be used. Either 'tweet' to use a public tweet,
    or 'profile' to use the profile bio.
    :param value: Value that should be used to get the data. Either the Tweet id if the type is 'tweet', or the
    username of the profile if type is 'profile'.
    :return The URL to which it is possible to get the data that can be used to verify the profile.
    :raise Exception if anything is wrong during the process. This can happen if:
            1. The tweet or profile bio does not contain any valid URL that link to a valid signature data object
            2. The provided signature is not valid
            3. The provided address is not linked to the provided public key
    """

    if type not in TYPES:
        raise Exception(f"Invalid type provided: {type}")

    # Get the URLs to check inside the tweet or the bio
    urls = []
    if type == TYPE_TWEET:
        urls = get_urls_from_tweet(value)
    elif type == TYPE_PROFILE:
        urls = get_urls_from_bio(value)

    if len(urls) == 0:
        raise Exception(f"No URL found inside {type}")

    # Find the signature following the URLs
    valid_url = ''
    data = None
    for url in urls:
        result = get_signature_from_url(url)
        if result is not None:
            valid_url = url
            data = result
            break

    if data is None:
        raise Exception(f"No valid signature data found inside {type}")

    # Verify the signature
    signature_valid = verify_signature(data)
    if not signature_valid:
        raise Exception("Invalid signature")

    # Verify the address
    address_valid = verify_address(data)
    if not address_valid:
        raise Exception("Invalid address")

    return valid_url


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
