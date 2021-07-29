use obi::{OBISchema, OBIEncode, OBIDecode};
use owasm::{execute_entry_point, oei, ext, prepare_entry_point};

const DATA_SOURCE_TWITTER: i64 = 49;
const DATA_SOURCE_GITHUB: i64 = 68;
const DATA_SOURCE_DISCORD: i64 = 80;

/// CallData contains the data that must be sent when calling this script
#[derive(OBIEncode, OBIDecode, OBISchema, Debug)]
struct CallData {
    application: String,
    call_data: String,
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
    } else if application == "github" {
        return DATA_SOURCE_GITHUB;
    } else if application == "discord" {
        return DATA_SOURCE_DISCORD;
    }

    panic!("Invalid application type")
}

#[no_mangle]
fn prepare_impl(input: CallData) {
    oei::ask_external_data(
        0,
        get_data_source(input.application),
        input.call_data.as_bytes(),
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

    #[test]
    fn test_obi_encode() {
        let input = CallData {
            application: "github".to_string(),
            call_data: "7B22757365726E616D65223A22526963636172646F4D222C22676973745F6964223A223732306530303732333930613930316262383065353966643630643766646564227D".to_string(),
        };

        let bytes = OBIEncode::try_to_vec(&input).unwrap();
        let hex_encoded = hex::encode(bytes);
        print!("Your call data is: {}", hex_encoded);
    }

    #[test]
    fn test_obi_encode_result() {
        let result = Result { value: "ricmontagnin".to_string(), signature: "655900061f9a306c685bbbf446a3cd02f469f995a2ead5dd64ae0ac0e90161b458a3a91dce3080fb935c98585ca2c1e93b1722ff52cdd5c59840f451438b8e2c".to_string() };
        let encode= OBIEncode::try_to_vec(&result).unwrap();
        println!("{}", base64::encode(encode));
    }

    #[test]
    fn test_obi_decode() {
        let result = "AAAAgGEwMGE3ZDViZDQ1ZTQyNjE1NjQ1ZmNhZWI0ZDgwMGFmMjI3MDRlNTQ5MzdhYjIzNWU1ZTUwYmViZDM4ZTg4Yjc2NWZkYjY5NmMyMjcxMmMwY2FiMTE3Njc1NmI2MzQ2Y2JjMTE0ODFjNTQ0ZDFmNzgyOGNiMjMzNjIwYzA2MTczAAAADHJpY21vbnRhZ25pbg==";
        let bytes = base64::decode(result).unwrap();
        let output: Result = OBIDecode::try_from_slice(&bytes).unwrap();
        print!("{:?}", output)
    }
}