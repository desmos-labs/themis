#!/usr/bin/env python3
import json
import sys
import requests
import re
from typing import Optional
import cryptography.hazmat.primitives.asymmetric.utils as crypto
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import hashes
import hashlib

ENDPOINT = "https://themis.morpheus.desmos.network/youtube"
HEADERS = {"Content-Type": "application/json"}

class CallData:
    """
    Contains the data that has been used to call the script
    """

    def __init__(self, user_id: str):
        self.user_id = user_id


class VerificationData:
    """
    Contains the data needed to verify the proof submitted by the user.
    """

    def __init__(self, address: str, pub_key: str, value: str, signature: str):
        self.address = address
        self.pub_key = pub_key
        self.signature = signature
        self.value = value


def get_urls_from_description(user: str) -> [str]:
    """
    Returns all the URLs that are found inside the description of the user having the given user id.
    :param user: Login name of the Youtube user.
    :return: List of URLs that are found inside the description
    """
    url = f"{ENDPOINT}/users/{user}"
    result = requests.request("GET", url, headers=HEADERS).json()
    return re.findall(r'(https?://[^\s]+)', result['description'])


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
    r = hashlib.new("ripemd160", s).digest()
    return data.address.upper() == r.hex().upper()


def check_values(values: dict) -> CallData:
    """
    Checks the validity of the given dictionary making sure it contains the proper data.
    :param values: Dictionary that should be checked.
    :return: A CallData instance.
    """
    if "user_id" not in values:
        raise Exception("Missing 'user_id' value")

    return CallData(values["user_id"])


def main(args: str):
    """
    Gets the signature data from Youtube reading it from a description.

    :param args Hex encoded JSON object containing the arguments to be used during the execution.
    In order to be valid, the encoded JSON object must contain a field named "user_id" that represents the
    Youtube user id of the user to be verified.

    The link must be public and its content must be a JSON object formed as follows:

    ```json
    {
      "address": "Hex encoded address of the signer",
      "pub_key": "Hex encoded public key that has been used to sign the value",
      "value": "Hex encoded value that has been signed",
      "signature": "Hex encoded Secp256k1 signature"
    }
    ```

    Example argument value:
    7B22757365726E616D65223A2272616E646F6D75736572227D

    This is the hex encoded representation of the following JSON object:

    ```json
    {
      "user_id":"randomuser"
    }
    ```

    :param args: JSON encoded parameters used during the execution.
    :return The signed value and the signature as a single comma separated string.
    :raise Exception if anything is wrong during the process. This can happen if:
            1. The user description does not contain any valid URL that link to a valid signature data object
            2. The provided signature is not valid
            3. The provided address is not linked to the provided public key
    """

    decoded = bytes.fromhex(args)
    json_obj = json.loads(decoded)
    call_data = check_values(json_obj)

    # Get the URLs to check from the user description
    urls = get_urls_from_description(call_data.user_id)
    if len(urls) == 0:
        raise Exception(f"No URL found inside {call_data.user_id} description")

    # Find the signature following the URLs
    data = None
    for url in urls:
        result = get_signature_from_url(url)
        if result is not None:
            data = result
            break

    if data is None:
        raise Exception(f"No valid signature data found inside {call_data.user_id} description")

    # Verify the signature
    signature_valid = verify_signature(data)
    if not signature_valid:
        raise Exception("Invalid signature")

    # Verify the address
    address_valid = verify_address(data)
    if not address_valid:
        raise Exception("Invalid address")

    return f"{data.value},{data.signature},{call_data.user_id}"


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)