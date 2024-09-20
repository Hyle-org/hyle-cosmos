pragma circom 2.1.6;

template HyleOutput(v,init_size,init_value,next_state_size,idy_size,idy_value,tx_size,tx_value,idx,payload_size,payload_value,is_success) {
 
   signal input next_value[next_state_size];

   signal output version;
   signal output initial_state[init_size];
   signal output next_state[next_state_size];
   signal output identity[idy_size];
   signal output tx_hash[tx_size];
   signal output index;
   signal output payloads[payload_size];
   signal output success;
   signal output length;

   version <== v;
   initial_state <== init_value;
   next_state <== next_value;
   if (idy_size > 0)
      identity <== idy_value;
   tx_hash <== tx_value;
   index <== idx;
   payloads <== payload_value;
   success <== is_success;

   length <== init_size + next_state_size + idy_size + tx_size + payload_size + 8;
}