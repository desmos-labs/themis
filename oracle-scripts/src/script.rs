use obi::{OBISchema, OBIEncode, OBIDecode};
use owasm::{execute_entry_point, oei, ext, prepare_entry_point};
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
    url: String,
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

/// Verifies the validity of the provided data, making sure that at least min_count of them are
/// have returned a valid URL.
///
/// If the data is valid, returns an Output instance that contains the URL that can be used
/// to verify the data provided by the user (which is returned by the data sources).
fn verify_validity(data: impl Iterator<Item=String>, min_count: i64) -> Output {
    let results = data.collect::<Vec<_>>();
    Output {
        valid: results.len() >= usize::try_from(min_count).unwrap(),
        url: results.get(0).unwrap().to_string(),
    }
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
    fn test_hex_byte_encoding() {
        let input = Input {
            validation_type: "tweet".to_string(),
            value: "1368883070590476292".to_string(),
        };

        let bytes = OBIEncode::try_to_vec(&input).unwrap();
        let hex_encoded = hex::encode(bytes);
        print!("{}", hex_encoded);
    }

    #[test]
    fn test_verify_validity() {
        let items = iter::repeat("https://t.co/bLokglOAel".to_string()).take(10);

        // Valid count
        let result = verify_validity(items.clone(), 5);
        assert_eq!(true, result.valid);
        assert_eq!("https://t.co/bLokglOAel".to_string(), result.url);

        // Invalid count
        let result = verify_validity(items.clone(), 15);
        assert_eq!(false, result.valid);
        assert_eq!("https://t.co/bLokglOAel".to_string(), result.url)
    }
}