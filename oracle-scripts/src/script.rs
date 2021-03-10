use obi::{OBISchema, OBIEncode, OBIDecode};
use owasm::{execute_entry_point, oei, ext, prepare_entry_point};
use serde::{Serialize, Deserialize};
use bitcoin_hashes::{sha256, Hash, hash160};
use secp256k1::{PublicKey, Secp256k1, Message, Signature};
use bech32::{self, ToBase32};
use std::convert::TryFrom;

const DESMOS_THEMIS_DS: i64 = 49;

#[derive(OBIEncode, OBIDecode, OBISchema)]
struct Input {
    validation_type: String,
    value: String,
}

#[derive(OBIEncode, OBISchema)]
struct Output {
    valid: bool,
}

#[no_mangle]
fn prepare_impl(input: Input) {
    // Get the data using the Desmos Themis data source
    oei::ask_external_data(
        0,
        DESMOS_THEMIS_DS,
        format!("{} {}", input.validation_type, input.value).as_bytes(),
    );
}

#[derive(Serialize, Deserialize)]
struct DataResult {
    address: String,
    pub_key: String,
    signature: String,
    value: String,
}

impl DataResult {
    fn is_valid(&self) -> bool {
        // Parse the public key
        let public_key = &self.pub_key.parse::<PublicKey>().unwrap();

        // Get the SHA-256 of the message
        let hash = sha256::Hash::hash(&self.value.as_bytes());
        let message = Message::from_slice(&hash).unwrap();

        // Parse the signature based on whether it is in a compact form (128 hex characters) or not
        let signature = match self.signature.len() == 128 {
            true => {
                // If the HEX encoding is long 128 characters, it means in a compact form and the
                // signature is actually composed on 128/2 = 64 characters
                let bytes = hex::decode(&self.signature).unwrap();
                Signature::from_compact(&bytes).unwrap()
            }
            false => {
                // If the HEX encoding is not long 128 characters, it means the signature is
                // encoded as extended
                self.signature.parse::<Signature>().unwrap()
            }
        };

        // Verify the signature
        let secp = Secp256k1::new();
        let is_signature_valid = secp.verify(&message, &signature, &public_key).is_ok();

        // Compute the address
        let ripemd160_hash = hash160::Hash::hash(&public_key.serialize());
        let (_, data) = bech32::decode(&self.address).unwrap();
        let is_address_correct = data.eq(&ripemd160_hash.to_vec().to_base32());

        return is_signature_valid && is_address_correct;
    }
}

/// Verifies the validity of the provided data, making sure that at least min_count of them are
/// effectively valid.
///
/// In order to be considered valid, each item present inside the iterator should:
/// - be a JSON object that can be parsed into a DataResult instance; and
/// - the parsed DataResult should contain valid data (the signature is correct and the
///   public key is associated with the given address)
///
/// After all the checks, it returns an Output object
fn verify_validity(data: impl Iterator<Item=String>, min_count: i64) -> Output {
    let result_valid = data
        .map(|i| {
            let value: DataResult = serde_json::from_str(&i).unwrap();
            value
        })
        .map(|data| data.is_valid())
        .count();

    Output { valid: result_valid >= usize::try_from(min_count).unwrap() }
}


#[no_mangle]
fn execute_impl(_input: Input) -> Output {
    // To make sure the data are valid we have to:
    // 1. get the results from the data sources
    // 2. map the results reading the returned value as DataResult objects
    // 3. make sure that all data results are valid (none is invalid)
    let result_valid = ext::load_input_raw(0);
    verify_validity(result_valid, oei::get_min_count())
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);

#[cfg(test)]
mod tests {
    use super::*;
    use hex;
    use std::iter;

    #[test]
    fn hex_byte_encoding() {
        let input = Input {
            validation_type: "tweet".to_string(),
            value: "1368883070590476292".to_string(),
        };

        let bytes = OBIEncode::try_to_vec(&input).unwrap();
        let hex_encoded = hex::encode(bytes);
        print!("{}", hex_encoded);
    }

    #[test]
    fn signature_check() {
        // Make sure that correct data are valid:
        // - the address is the one associated with the given pub key
        // - the signature is valid considering the public key and value provided
        let data = DataResult {
            address: "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu".to_string(),
            pub_key: "033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d".to_string(),
            signature: "a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173".to_string(),
            value: "ricmontagnin".to_string(),
        };
        assert_eq!(true, data.is_valid())
    }

    #[test]
    fn execute_impl() {
        let data = DataResult {
            address: "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu".to_string(),
            pub_key: "033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d".to_string(),
            signature: "a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173".to_string(),
            value: "ricmontagnin".to_string(),
        };
        let string = serde_json::to_string(&data).unwrap();
        let items = iter::repeat(string).take(10);

        // Valid count
        let result = verify_validity(items.clone(), 5);
        assert_eq!(true, result.valid);

        // Invalid count
        let result = verify_validity(items.clone(), 15);
        assert_eq!(false, result.valid)
    }
}