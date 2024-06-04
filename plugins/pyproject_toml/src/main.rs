use prost::Message;
use std::io;
use std::io::prelude::*;

pub mod plugin {
    include!(concat!(env!("OUT_DIR"), "/plugin.rs"));
}

pub fn deserialize_get_version_request(
    buf: &[u8],
) -> Result<plugin::GetVersionRequest, prost::DecodeError> {
    // return plugin::GetVersionRequest::decode(&mut Cursor::new(buf));

    plugin::GetVersionRequest::decode(buf)
}

pub fn serialize_codegen_response(resp: &plugin::GetVersionResponse) -> Vec<u8> {
    let mut buf = Vec::with_capacity(resp.encoded_len());

    resp.encode(&mut buf).unwrap();
    buf
}

pub fn create_get_version_response(req: plugin::GetVersionRequest) -> plugin::GetVersionResponse {
    let file = req.inputs;
    if let Some(file) = file {
        let mut file = std::fs::File::open(file.path).unwrap();
        let mut contents = String::new();
        file.read_to_string(&mut contents).unwrap();
        return plugin::GetVersionResponse {
            version: contents.trim().to_string(),
        };
    }

    plugin::GetVersionResponse {
        version: "".to_string(),
    }
}

#[allow(dead_code)]
fn main() -> Result<(), prost::DecodeError> {
    let stdin = io::stdin();
    let mut stdin = stdin.lock();
    let buffer = stdin.fill_buf().unwrap();

    let req = match deserialize_get_version_request(buffer) {
        Ok(request_deserialized_result) => request_deserialized_result,
        Err(_e) => std::process::exit(1),
    };

    let resp = create_get_version_response(req);
    let out = serialize_codegen_response(&resp);

    match io::stdout().write_all(&out) {
        Ok(result) => result,
        Err(_e) => std::process::exit(1),
    };

    Ok(())
}
