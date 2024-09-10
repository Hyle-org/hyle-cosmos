use hex::FromHex;
use risc0_zkvm::sha::Digest;
use serde_json;
use std::env;

use hyle_contract::HyleOutput;

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() != 3 {
        eprintln!("Usage: {} <image_id> <receipt_path>", args[0]);
        std::process::exit(1);
    }

    // Image ID is the hexademical representation of the method ID, without leading prefix.
    let image_id = &args[1];

    // Parse the proof from file
    let receipt_path = &args[2];
    let receipt_content =
        std::fs::read_to_string(receipt_path).expect("Failed to read receipt file");
    let receipt: risc0_zkvm::Receipt =
        serde_json::from_str(&receipt_content).expect("Failed to parse receipt file");
    let image_bytes = Digest::from_hex(image_id).expect("Invalid image ID hex");

    // perform verification
    receipt.verify(image_bytes).expect("Verification failed");

    // Outputs to stdout for the caller to read.
    let output: HyleOutput<()> = receipt
        .journal
        .decode()
        .expect("Failed to decode receipt journal");
    println!(
        "{}",
        serde_json::to_string(&output).expect("Failed to serialize output")
    );
}
