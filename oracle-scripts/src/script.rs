use obi::{OBISchema, OBIEncode, OBIDecode};
use owasm::{execute_entry_point, oei, ext, prepare_entry_point};
use std::convert::TryFrom;

const APP_TWITTER: &'static str = "twitter";
const DATA_SOURCE_TWITTER: i64 = 49;

/// VerificationData contains the data used to verify the ownership of an application account
#[derive(OBIEncode, OBIDecode, OBISchema, Debug)]
struct VerificationData {
    method: String,
    value: String,
}

/// CallData contains the data that must be sent when calling this script
#[derive(OBIEncode, OBIDecode, OBISchema, Debug)]
struct CallData {
    application: String,
    verification_data: VerificationData,
}

/// SignatureData contains the data of a signature after it has been verified
#[derive(OBIEncode, OBIDecode, OBISchema, Debug)]
struct SignatureData {
    signature: String,
    value: String,
}

/// Result contains the data of the execution result
#[derive(OBIEncode, OBIDecode, OBISchema, Debug)]
struct Result {
    valid: bool,
    signature_data: SignatureData,
}

#[no_mangle]
fn prepare_impl(input: CallData) {
    let mut data_source_id = 0;
    if input.application == APP_TWITTER {
        data_source_id = DATA_SOURCE_TWITTER
    }

    oei::ask_external_data(0, data_source_id, format!(
        "{} {}",
        input.verification_data.method,
        input.verification_data.value,
    ).as_bytes());
}

/// Verifies the validity of the provided data, making sure that at least min_count of them are
/// have returned a valid URL.
///
/// If the data is valid, returns an Output instance that contains the URL that can be used
/// to verify the data provided by the user (which is returned by the data sources).
fn verify_validity(data: impl Iterator<Item=String>, min_count: i64) -> Result {
    let results = data.collect::<Vec<_>>();

    // Read the returned data
    let data = results.get(0).unwrap().to_string();
    let parts = data.split(",").collect::<Vec<&str>>();

    Result {
        valid: results.len() >= usize::try_from(min_count).unwrap(),
        signature_data: SignatureData {
            value: parts[0].to_string(),
            signature: parts[1].to_string().replace("\n", ""),
        },
    }
}

#[no_mangle]
fn execute_impl(_input: CallData) -> Result {
    if _input.application != APP_TWITTER {
        return Result {
            valid: false,
            signature_data: SignatureData { value: "".to_string(), signature: "".to_string() },
        };
    }

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
    fn test_obi_encode() {
        let input = CallData {
            application: "twitter".to_string(),
            verification_data: VerificationData {
                method: "tweet".to_string(),
                value: "1392033585675317252".to_string(),
            },
        };

        let bytes = OBIEncode::try_to_vec(&input).unwrap();
        let hex_encoded = hex::encode(bytes);
        print!("Your call data is: {}", hex_encoded);
    }

    #[test]
    fn test_obi_decode() {
        let result = "AQAAAIBhMDBhN2Q1YmQ0NWU0MjYxNTY0NWZjYWViNGQ4MDBhZjIyNzA0ZTU0OTM3YWIyMzVlNWU1MGJlYmQzOGU4OGI3NjVmZGI2OTZjMjI3MTJjMGNhYjExNzY3NTZiNjM0NmNiYzExNDgxYzU0NGQxZjc4MjhjYjIzMzYyMGMwNjE3MwAAAAxyaWNtb250YWduaW4=";
        let bytes = base64::decode(result).unwrap();
        let output: Result = OBIDecode::try_from_slice(&bytes).unwrap();
        print!("{:?}", output)
    }

    #[test]
    fn test_verify_validity() {
        let items = iter::repeat("https://t.co/bLokglOAel".to_string()).take(10);

        // Valid count
        let result = verify_validity(items.clone(), 5);
        assert_eq!(true, result.valid);
        assert_eq!("ricmontagnin".to_string(), result.signature_data.value);
        assert_eq!(
            "a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173".to_string(),
            result.signature_data.signature,
        );

        // Invalid count
        let result = verify_validity(items.clone(), 15);
        assert_eq!(false, result.valid);
        assert_eq!("ricmontagnin".to_string(), result.signature_data.value);
        assert_eq!(
            "a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173".to_string(),
            result.signature_data.signature,
        );
    }
}