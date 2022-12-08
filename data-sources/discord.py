#!/usr/bin/env python3
import json
import sys
import urllib.parse
import requests
from typing import Optional
import cryptography.hazmat.primitives.asymmetric.utils as crypto
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import hashes
from ripemd.ripemd160 import ripemd160
import hashlib

ENDPOINT = "https://themis.morpheus.desmos.network/discord"
HEADERS = {"Content-Type": "application/json"}


class CallData:
    """
    Contains the data that has been used to call the script
    """

    def __init__(self, username: str):
        self.username = username


class VerificationData:
    """
    Contains the data needed to verify the proof submitted by the user.
    """

    def __init__(self, address: str, pub_key: str, value: str, signature: str):
        self.address = address
        self.pub_key = pub_key
        self.signature = signature
        self.value = value


def get_user_data(data: CallData) -> Optional[VerificationData]:
    """
    Tries getting the verification data for the user having the given Discord username.
    :param data: Data used to get the VerificationData
    :return: An OptionalData object if the call was successful, or None if it errored somehow.
    """
    try:
        url_encoded_username = urllib.parse.quote(data.username)
        result = requests.request("GET", f"{ENDPOINT}/{url_encoded_username}", headers=HEADERS).json()
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
    if "username" not in values:
        raise Exception("Missing 'username' value")

    return CallData(values["username"])


def main(args: str):
    """
    Gets the signature data from Discord, after the user has provided it through the Hephaestus Discord bot.

    :param args Hex encoded JSON object containing the arguments to be used during the execution.
    In order to be valid, the encoded JSON object must contain one field named "username" that represents the Discord
    username of the account to be connected.

    Example argument value:
    7B22757365726E616D65223A22526963636172646F204D6F6E7461676E696E2335343134227D

    This is the hex encoded representation of the following JSON object:

    ```json
    {
      "username":"Riccardo Montagnin#5414"
    }
    ```

    :param args: JSON encoded parameters used during the execution.
    :return The signed value and the signature as a single comma separated string.
    :raise Exception if anything is wrong during the process. This can happen if:
            1. The Discord user has not started the connection using the Hephaestus bot
            2. The provided signature is not valid
            3. The provided address is not linked to the provided public key
    """

    decoded = bytes.fromhex(args)
    json_obj = json.loads(decoded)
    call_data = check_values(json_obj)

    result = get_user_data(call_data)
    if result is None:
        raise Exception(f"No valid signature data found for user with username {call_data.username}")

    # Verify the signature
    signature_valid = verify_signature(result)
    if not signature_valid:
        raise Exception("Invalid signature")

    # Verify the address
    address_valid = verify_address(result)
    if not address_valid:
        raise Exception("Invalid address")

    return f"{result.value},{result.signature},{call_data.username}"


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
