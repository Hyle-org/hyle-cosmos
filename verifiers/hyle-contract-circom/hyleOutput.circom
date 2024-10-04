pragma circom 2.1.6;

template HyleOutput(init_state_size,next_state_size,identity_size,tx_size,payload_size) {
   //input signal
   signal input version_input;
   signal input initial_state_input[init_state_size];
   signal input next_state_input[next_state_size];
   signal input identity_input[identity_size];
   signal input tx_hash_input[tx_size];
   signal input index_input;
   signal input payloads_input[payload_size];
   signal input success_input;

   //output signal
   signal output version_output;
   signal output initial_state_output[init_state_size];
   signal output next_state_output[next_state_size];
   signal output identity_output[identity_size];
   signal output tx_hash_output[tx_size];
   signal output index_output;
   signal output payloads_output[payload_size];
   signal output success_output;
   signal output length_output;

   version_output <== version_input;
   initial_state_output <== initial_state_input;
   next_state_output <== next_state_input;
   if (identity_size > 0)
      identity_output <== identity_input;
   tx_hash_output <== tx_hash_input;
   index_output <== index_input;
   payloads_output <== payloads_input;
   success_output <== success_input;

   //length calculation
   length_output <== init_state_size + next_state_size + identity_size + tx_size + payload_size + 8;
}