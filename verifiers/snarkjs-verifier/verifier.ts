import * as fs from "fs";
import { parseArgs } from "util";
import { spawn } from "child_process";

interface HyleOutput {
    version: number;
    initial_state: number[];
    next_state: number[];
    identity: string;
    tx_hash: number[];
    index: number;
    payloads: number[];
    success: boolean;
}

function parseString(vector: number[]): string {
    let length = vector.length;
    let resp = "";
    for (var i = 0; i < length; i++) resp += String.fromCharCode(vector[i], 16);
    return resp;
}

function deserializePublicInputs<T>(publicInput: string[]): HyleOutput {
    const nbValueToExtract = parseInt(publicInput[publicInput.length - 1], 10);
    const extractedData = publicInput.slice(publicInput.length-nbValueToExtract-1);
    //version
    const version = parseInt(extractedData.shift() as string);
    //initial_state
    const initial_state_len = parseInt(extractedData.shift() as string);
    const initial_state: number[] = [];
    for (let i = 0; i < initial_state_len; i++) {
        initial_state.push(parseInt(extractedData.shift() as string));
    }
    //next_state
    const next_state_len = parseInt(extractedData.shift() as string);
    const next_state: number[] = [];
    for (let i = 0; i < next_state_len; i++) {
        next_state.push(parseInt(extractedData.shift() as string));
    }
    //idenity
    const identity_len = parseInt(extractedData.shift() as string);
    const identity_array: number[] = [];
    for (let i = 0; i < identity_len; i += 1) {
        identity_array.push(parseInt(extractedData.shift() as string));
    }
    const identity = identity_len > 0 ? parseString(identity_array) : "";
    //tx_hash
    const tx_hash_len = parseInt(extractedData.shift() as string);
    const tx_hash: number[] = [];
    for (let i = 0; i < tx_hash_len; i += 1) {
        tx_hash.push(parseInt(extractedData.shift() as string));
    }
    //index
    const index = parseInt(extractedData.shift() as string);
    //payloads
    const payloads_len = parseInt(extractedData.shift() as string);
    const payloads: number[] = [];
    for (let i = 0; i < payloads_len; i += 1) {
        payloads.push(parseInt(extractedData.shift() as string));
    }
    //success
    const success = parseInt(extractedData.shift() as string) === 1;
    // We don't parse the rest, which correspond to programOutputs
    return {
        version,
        initial_state,
        next_state,
        identity,
        tx_hash,
        index,
        payloads,
        success,
    };
}

function runCommand(command: string, args: string[]): Promise<string> {
    return new Promise((resolve, reject) => {
        const process = spawn(command, args);

        let stdoutData = "";
        let stderrData = "";

        process.stdout.on("data", (data) => {
            stdoutData += data.toString()
        });

        process.stderr.on("data", (data) => {
            stderrData += data.toString()
        });

        process.on("close", (code) => {
            if (code === 0) {
                resolve(stdoutData.trim());
            } else {
                reject(new Error(`Process exited with code ${code}`));
            }
        });
    });
}

async function main() {
    const { values, positionals } = parseArgs({
        args: process.argv,
        options: {
            vKeyPath: {
                type: "string",
            },
            proofPath: {
                type: "string",
            },
            publicInput: {
                type: "string",
            }
        },
        strict: true,
        allowPositionals: true,
    });

    const command = "bash";

    const argsVerification = ["-c", `snarkjs groth16 verify ${values.vKeyPath} ${values.publicInput} ${values.proofPath}`];
    const result: string = await runCommand(command, argsVerification);

    // Proof is considered valid
    if (result.includes('OK')) {
        const jsonData = JSON.parse(fs.readFileSync(values.publicInput));
        const hyleOutput = deserializePublicInputs(jsonData);

        var stringified_output = JSON.stringify(hyleOutput);

        process.stdout.write(stringified_output);
        process.exit(0);
    } else {
        throw Error("Invalid proof");
    }
}

try {
    await main();
} catch (e) {
    console.error(e);
    process.exit(1);
}