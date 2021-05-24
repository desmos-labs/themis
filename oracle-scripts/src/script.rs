use obi::{OBISchema, OBIEncode, OBIDecode};
use owasm::{execute_entry_point, oei, ext, prepare_entry_point};

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

/// Result contains the data of the execution result
#[derive(OBIEncode, OBIDecode, OBISchema, Debug)]
struct Result {
    signature: String,
    value: String,
}

/// Returns the correct data source based on the given input
fn get_data_source(application: String) -> i64 {
    if application == "twitter" {
        return DATA_SOURCE_TWITTER;
    }

    panic!("Invalid application type")
}

#[no_mangle]
fn prepare_impl(input: CallData) {
    oei::ask_external_data(
        0,
        get_data_source(input.application),
        format!(
            "{} {}",
            input.verification_data.method,
            input.verification_data.value,
        ).as_bytes(),
    );
}

#[no_mangle]
fn execute_impl(_input: CallData) -> Result {
    // To make sure the data are valid we have to:
    // 1. get the results from the data sources
    // 2. map the results reading the returned value as DataResult objects
    // 3. make sure that all data results are valid (none is invalid)
    let valid_results = ext::load_input_raw(0).collect::<Vec<_>>();

    // Read the returned data
    let data = valid_results.get(0).unwrap().to_string();
    let parts = data.split(",").collect::<Vec<&str>>();

    Result {
        value: parts[0].to_string(),
        signature: parts[1].to_string().replace("\n", ""),
    }
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
            application: "github".to_string(),
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
        let result = "AAAAgGEwMGE3ZDViZDQ1ZTQyNjE1NjQ1ZmNhZWI0ZDgwMGFmMjI3MDRlNTQ5MzdhYjIzNWU1ZTUwYmViZDM4ZTg4Yjc2NWZkYjY5NmMyMjcxMmMwY2FiMTE3Njc1NmI2MzQ2Y2JjMTE0ODFjNTQ0ZDFmNzgyOGNiMjMzNjIwYzA2MTczAAAADHJpY21vbnRhZ25pbg==";
        let bytes = base64::decode(result).unwrap();
        let output: Result = OBIDecode::try_from_slice(&bytes).unwrap();
        print!("{:?}", output)
    }

    #[test]
    fn test_verify_validity() {
        let items = iter::repeat("https://t.co/bLokglOAel".to_string()).take(10);

        let result = verify_validity(items.clone());
        assert_eq!("ricmontagnin".to_string(), result.value);
        assert_eq!(
            "a00a7d5bd45e42615645fcaeb4d800af22704e54937ab235e5e50bebd38e88b765fdb696c22712c0cab1176756b6346cbc11481c544d1f7828cb233620c06173".to_string(),
            result.signature,
        );
    }
}