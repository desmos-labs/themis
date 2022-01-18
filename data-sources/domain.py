#!/usr/bin/env python3
import json
import sys
import urllib.parse
import requests
from typing import Optional
import cryptography.hazmat.primitives.asymmetric.utils as crypto
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import hashes
import hashlib

ENDPOINT = "https://themis.morpheus.desmos.network/nslookup"
HEADERS = {"Content-Type": "application/json"}


class CallData:
    """
    Contains the data that has been used to call the script
    """

    def __init__(self, domain: str):
        self.domain = domain


class VerificationData:
    """
    Contains the data needed to verify the proof submitted by the user.
    """

    def __init__(self, address: str, pub_key: str, value: str, signature: str):
        self.address = address
        self.pub_key = pub_key
        self.signature = signature
        self.value = value


def validate_json(json: dict) -> bool:
    """
    Tells whether or not the given JSON is a valid signature JSON object.
    :param json: JSON object to be checked.
    :return: True if the provided JSON has a valid signature schema, or False otherwise.
    """

    return all(key in json for key in ['value', 'pub_key', 'signature', 'address'])


def try_reading_json(json_value: dict) -> Optional[VerificationData]:
    """
    Tries reading the given param as a JSON object.
    :param json_value: Value to be checked for validity.
    :return: A VerificationData instance if the value is a valid JSON object containing the needed
    fields, or None otherwise.
    """

    if validate_json(json_value):
        return VerificationData(
            json_value['address'],
            json_value['pub_key'],
            json_value['value'],
            json_value['signature'],
        )

    return None


def get_data_from_txt_record(record: str) -> Optional[VerificationData]:
    """
    Tries getting a VerificationData from the specified record value.
    Before failing and returning None, it tries two different ways:
    1. Reading the value as a JSON object containing the required fields.
    2. Reading the value as a URL that contains a JSON object with the required fields.

    :param record: Value of the TXT record that should be used.
    :return: Either a VerificationData upon a successful check, or None if the record value does not point
    to a valid JSON object anyway.
    """
    try:
        value = try_reading_json(json.loads(record))
        if value is not None:
            return value

    except ValueError:
        try:
            content = requests.request("GET", record, headers=HEADERS).json()
            value = try_reading_json(content)
            if value is not None:
                return value

        except ValueError:
            return None

    return None


def get_user_data(data: CallData) -> Optional[VerificationData]:
    """
    Tries getting the verification data by reading the TXT records of a specific domain name.
    :param data: Data used to get the VerificationData
    :return: An OptionalData object if the call was successful, or None if it errored somehow.
    """
    try:
        url_encoded_domain = urllib.parse.quote(data.domain)
        response = requests.request("GET", f"{ENDPOINT}/{url_encoded_domain}", headers=HEADERS)
        if response.status_code != 200:
            return None

        # For each record we need to check its text value and either try to parse the JSON or get it from the link
        for record in response.json()['txt']:
            value = get_data_from_txt_record(record['text'])
            if value is not None:
                return value

    except ValueError as err:
        print(err)
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
    if "domain" not in values:
        raise Exception("Missing 'domain' value")

    return CallData(values["domain"])


def main(args: str):
    """
    Gets the signature data from a domain TXT records, after the user has updated their values.
    In order to be a valid entry used during the verification process, a TXT record can either value it's value set to:
    a. a JSON object containing the required fields ('value', 'pub_key', 'signature', 'address').
    b. a URL pointing to a page containing a valid JSON object with the required fields ('value', 'pub_key',
       'signature', 'address').

    :param args Hex encoded JSON object containing the arguments to be used during the execution.
    In order to be valid, the encoded JSON object must contain one field named "domain" that represents the domain
    name which TXT records should be used as the verification method.

    Example argument value:
    7B22646F6D61696E223A22666F72626F6C652E636F6D227D

    This is the hex encoded representation of the following JSON object:

    ```json
    {
      "domain":"forbole.com"
    }
    ```

    :param args: JSON encoded parameters used during the execution.
    :return The signed value and the signature as a single comma separated string.
    :raise Exception if anything is wrong during the process. This can happen if:
            1. The domain name has not updated its TXT records yet.
            2. The provided signature is not valid
            3. The provided address is not linked to the provided public key
    """

    decoded = bytes.fromhex(args)
    json_obj = json.loads(decoded)
    call_data = check_values(json_obj)

    result = get_user_data(call_data)
    if result is None:
        raise Exception(f"No valid signature data found for domain {call_data.domain}")

    # Verify the signature
    signature_valid = verify_signature(result)
    if not signature_valid:
        raise Exception("Invalid signature")

    # Verify the address
    address_valid = verify_address(result)
    if not address_valid:
        raise Exception("Invalid address")

    return f"{result.value},{result.signature},{call_data.domain}"


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
