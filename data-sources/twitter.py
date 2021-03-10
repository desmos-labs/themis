#!/usr/bin/env python3

import sys
import requests
import re
import cryptography.hazmat.primitives.asymmetric.utils as crypto
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import hashes
import hashlib
import bech32

TYPE_TWEET = "tweet"
TYPE_PROFILE = "profile"
TYPES = [TYPE_TWEET, TYPE_PROFILE]

ENDPOINT = "https://themis.morpheus.desmos.network/twitter"
HEADERS = {"Content-Type": "application/json"}


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


def get_signature_from_url(url: str) -> dict:
    """
    Tries getting the signature object linked to the given URL.
    :param url: URL that should contain the signature object.
    :return: A dictionary containing 'valid' to tell whether the search was valid, and an optional 'data' containing
    the signature object.
    """
    try:
        result = requests.request("GET", url, headers=HEADERS).json()
        if validate_json(result):
            return {'valid': True, 'data': result}
        else:
            return {'valid': False}
    except ValueError:
        return {'valid': False}


def verify_signature(pubkey: str, signature: str, value: str) -> bool:
    """
    Verifies the signature using the given pubkey and value.
    :param pubkey: HEX encoded Secp256k1 public key that should be used to verify the signature.
    :param signature: HEX encoded signature of the value.
    :param value: Value that has been signed.
    :return True if the signature is valid, False otherwise
    """
    if len(signature) != 128:
        return False

    try:
        # Create signature for dss signature
        (r, s) = int(signature[:64], 16), int(signature[64:], 16)
        sig = crypto.encode_dss_signature(r, s)

        # Create public key instance
        public_key = ec.EllipticCurvePublicKey.from_encoded_point(ec.SECP256K1(), bytes.fromhex(pubkey))

        # Verify the signature
        public_key.verify(sig, str.encode(value), ec.ECDSA(hashes.SHA256()))
        return True
    except Exception:
        return False


def verify_address(address: str, pubkey: str) -> bool:
    """
    Verifies that the given address is the one associated with the provided HEX encoded compact public key.
    :param address: Bech32 encoded address that should be checked.
    :param pubkey: HEX encoded public key in the compact form
    :return: True if the given address is associated to the given public key, False otherwise
    """
    s = hashlib.new("sha256", bytes.fromhex(pubkey)).digest()
    r = hashlib.new("ripemd160", s).digest()
    five_bit_r = bech32.convertbits(r, 8, 5)
    address_bytes = bech32.bech32_decode(address)
    return five_bit_r == address_bytes


def main(type, value):
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
    data = {}
    for url in urls:
        result = get_signature_from_url(url)
        if result['valid']:
            valid_url = url
            data = result['data']
            break

    if not data:
        raise Exception(f"No valid signature data found inside {type}")

    # Unpack the values
    address, signature, pub_key, value = data['address'], data['signature'], data['pub_key'], data['value']

    # Verify the signature
    signature_valid = verify_signature(pub_key, signature, value)
    if not signature_valid:
        raise Exception("Invalid signature")

    # Verify the address
    address_valid = verify_address(address, pub_key)
    if not address_valid:
        raise Exception("Invalid address")

    return valid_url


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
