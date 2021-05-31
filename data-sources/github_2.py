
import json
import sys
import requests
from typing import Optional
import cryptography.hazmat.primitives.asymmetric.utils as crypto
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives import hashes

import hashlib

HEADERS = {"Content-Type": "application/json"}


class CallData:
    """
    Contains the data that has been used to call the script
    """

    def __init__(self, username: str, gist_id: str):
        self.username = username
        self.gist_id = gist_id


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


def get_data_from_gist(data: CallData) -> Optional[VerificationData]:
    """
    Tries getting the signature object from within the Gist having the given id.
    :param data: Object containing the data that has been used to call the script.
    :return: A VerificationData instance if no error is raised, or None otherwhise.
    """
    try:
        url = f"https://gist.githubusercontent.com/{data.username}/{data.gist_id}/raw/"
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


def check_values(values: dict) -> CallData:
    """
    Checks the validity of the given dictionary making sure it contains the proper data.
    :param values: Dictionary that should be checked.
    :return: A CallData instance.
    """
    if "username" not in values:
        raise Exception("Missing username value")

    if "gist_id" not in values:
        raise Exception("Missing gist_id value")

    return CallData(values["username"], values["gist_id"])


def main(args: str):
    """
    Gets the signature data from GitHub, using a public Gist.

    :param args: Hex encoded JSON object containing the arguments to be used during the execution.
    In order to be valid, the encoded JSON object must contain two fields:

    1. "username", which represents the GitHub name of the user that is verifying their identity;
    2. "gist_id", which is the ID of the Gist that the user has created.

    The "gist_id" must reference a **public** Gist that contains a JSON object formed as follows:

    ```json
    {
      "address": "Hex address os the signer",
      "pub_key": "Hex encoded public key that has been used to sign the value",
      "value": "Value that has been signed",
      "signature": "Hex encoded Secp256k1 signature"
    }
    ```

    Example argument value:
    7B22757365726E616D65223A22526963636172646F4D222C22676973745F6964223A223732306530303732333930613930316262383065353966643630643766646564227D

    This is the hex encoded representation of the following JSON object:

    ```json
    {
      "username":"RiccardoM",
      "gist_id":"720e0072390a901bb80e59fd60d7fded"
    }
    ```

    :return The signed value and the signature as a single comma separated string.
    :raise Exception if anything is wrong during the process. This can happen if:
            1. The Gist is not publicly reachable.
            2. The provided signature is not valid
            3. The provided address is not linked to the provided public key
    """

    decoded = bytes.fromhex(args)
    json_obj = json.loads(decoded)
    call_data = check_values(json_obj)

    # Get the URLs to check inside the tweet or the bio
    data = get_data_from_gist(call_data)
    if data is None:
        raise Exception(f"No valid signature data found for gist with id {call_data.gist_id}")

    if data.value != call_data.username:
        raise Exception("Signed value is different from the provided GitHub username")

    # Verify the signature
    signature_valid = verify_signature(data)
    if not signature_valid:
        raise Exception("Invalid signature")

    # Verify the address
    address_valid = verify_address(data)
    if not address_valid:
        raise Exception("Invalid address")

    return f"{data.value},{data.signature}"


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
