package keeper_test

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/hyle-org/hyle/x/zktx"
	"github.com/hyle-org/hyle/x/zktx/keeper/gnark"
	"github.com/stretchr/testify/require"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/secp256k1/ecdsa"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/hash/sha3"
	"github.com/consensys/gnark/std/math/bitslice"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/selector"

	nativesha3 "github.com/ethereum/go-ethereum/crypto"

	circuitecdsa "github.com/consensys/gnark/std/signature/ecdsa"
)

// GNARK circuit for ECDSA verification, this is implemented in emulated arithmetic so it's inefficient.
type ecdsaCircuit[T, S emulated.FieldParams] struct {
	gnark.HyleCircuit
	Sig circuitecdsa.Signature[S]
	Msg emulated.Element[S] `gnark:",public"`
	Pub circuitecdsa.PublicKey[T, S]
}

func (c *ecdsaCircuit[T, S]) Define(api frontend.API) error {
	// Verify address
	newHasher, err := sha3.NewLegacyKeccak256(api)
	if err != nil {
		return err
	}

	uapi, err := uints.New[uints.U64](api)
	if err != nil {
		return err
	}
	pubKeyBytes, err := pubKeyToBytes(api, &c.Pub)
	if err != nil {
		return err
	}
	newHasher.Write(pubKeyBytes)
	res := newHasher.Sum()

	hexChars := []frontend.Variable{
		[]byte("0")[0],
		[]byte("1")[0],
		[]byte("2")[0],
		[]byte("3")[0],
		[]byte("4")[0],
		[]byte("5")[0],
		[]byte("6")[0],
		[]byte("7")[0],
		[]byte("8")[0],
		[]byte("9")[0],
		[]byte("a")[0],
		[]byte("b")[0],
		[]byte("c")[0],
		[]byte("d")[0],
		[]byte("e")[0],
		[]byte("f")[0],
	}

	// no Reset() on this hasher so make a new one
	if newHasher, err = sha3.NewLegacyKeccak256(api); err != nil {
		return err
	}
	newHasher.Write(limbsToBytes(uapi, c.Sig.R.Limbs))
	newHasher.Write(limbsToBytes(uapi, c.Sig.S.Limbs))
	payloadHash := newHasher.Sum()
	if a, b := len(payloadHash), len(c.HyleCircuit.PayloadHash); a != b {
		return fmt.Errorf("payload hash length mismatch %d != %d", a, b)
	}
	for i := 0; i < len(payloadHash); i++ {
		uapi.ByteAssertEq(payloadHash[i], c.HyleCircuit.PayloadHash[i])
	}

	for i := 0; i < 20; i++ {
		// Not sure if there's a more efficient way to do this but it works - we need to compare the ASCII values.
		lower, upper := bitslice.Partition(api, res[i+12].Val, 4)
		lower = selector.Mux(api, lower, hexChars[:]...)
		upper = selector.Mux(api, upper, hexChars[:]...)
		uapi.ByteAssertEq(uapi.ByteValueOf(upper), c.Identity[i*2])
		uapi.ByteAssertEq(uapi.ByteValueOf(lower), c.Identity[i*2+1])
	}
	// Check that the next one is a dot, aka a name separator
	uapi.ByteAssertEq(c.Identity[40], uints.NewU8(46))

	c.Pub.Verify(api, sw_emulated.GetCurveParams[T](), &c.Msg, &c.Sig)
	return nil
}

// The following two functions directly pulled from https://github.com/Consensys/gnark/discussions/802
func pubKeyToBytes[T, S emulated.FieldParams](api frontend.API, pubKey *circuitecdsa.PublicKey[T, S]) ([]uints.U8, error) {
	xLimbs := pubKey.X.Limbs
	yLimbs := pubKey.Y.Limbs

	u64api, err := uints.New[uints.U64](api)
	if err != nil {
		return nil, err
	}

	result := limbsToBytes(u64api, xLimbs)
	return append(result, limbsToBytes(u64api, yLimbs)...), nil
}

func limbsToBytes(u64api *uints.BinaryField[uints.U64], limbs []frontend.Variable) []uints.U8 {
	result := make([]uints.U8, 0, len(limbs)*8)
	for i := range limbs {
		u64 := u64api.ValueOf(limbs[len(limbs)-1-i])
		result = append(result, u64api.UnpackMSB(u64)...)
	}
	return result
}

func sign(privKey *ecdsa.PrivateKey, msg []byte) ([]byte, error) {
	byteMsg := []byte(msg)
	md := sha256.New()
	sig, _ := privKey.Sign(byteMsg, md)

	// check that the signature is correct
	publicKey := privKey.PublicKey
	flag, _ := publicKey.Verify(sig, byteMsg, md)
	if !flag {
		return sig, fmt.Errorf("invalid signature")
	}
	return sig, nil
}

func makeECDSACircuit(privKey *ecdsa.PrivateKey, ethAddress string, sigBin, msg []byte) (*ecdsaCircuit[emulated.Secp256k1Fp, emulated.Secp256k1Fr], error) {
	// unmarshal signature
	var sig ecdsa.Signature
	sig.SetBytes(sigBin)
	r, s := new(big.Int), new(big.Int)
	r.SetBytes(sig.R[:32])
	s.SetBytes(sig.S[:32])

	// compute the hash of the message as an integer
	dataToHash := make([]byte, len(msg))
	copy(dataToHash[:], msg[:])
	md := sha256.New()
	md.Write(dataToHash[:])
	hramBin := md.Sum(nil)
	hash := ecdsa.HashToInt(hramBin)

	payloadHash, err := computePayloadHash(sigBin)
	if err != nil {
		return nil, err
	}

	return &ecdsaCircuit[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		HyleCircuit: gnark.HyleCircuit{
			Version:     1,
			InputLen:    1,
			Input:       []frontend.Variable{0},
			OutputLen:   1,
			Output:      []frontend.Variable{0},
			IdentityLen: len(ethAddress),
			Identity:    gnark.ToArray256([]byte(ethAddress)), // We expect only the origin as this is the "auth contract"
			TxHash:      gnark.ToArray64([]byte("TODO")),
			Success:     1,
			PayloadHash: gnark.ToArray32(payloadHash),
		},
		Sig: circuitecdsa.Signature[emulated.Secp256k1Fr]{
			R: emulated.ValueOf[emulated.Secp256k1Fr](r),
			S: emulated.ValueOf[emulated.Secp256k1Fr](s),
		},
		Msg: emulated.ValueOf[emulated.Secp256k1Fr](hash),
		Pub: circuitecdsa.PublicKey[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](privKey.PublicKey.A.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](privKey.PublicKey.A.Y),
		},
	}, nil
}

func generate_ecdsa_proof(circuit *ecdsaCircuit[emulated.Secp256k1Fp, emulated.Secp256k1Fr]) (gnark.Groth16Proof, error) {
	err := test.IsSolved(circuit, circuit, ecc.BN254.ScalarField())
	if err != nil {
		return gnark.Groth16Proof{}, err
	}

	// Witness first then compile as that modifies the circuit
	witness, err := frontend.NewWitness(circuit, ecc.BN254.ScalarField())
	if err != nil {
		return gnark.Groth16Proof{}, err
	}

	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		return gnark.Groth16Proof{}, err
	}

	// generating pk, vk
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return gnark.Groth16Proof{}, err
	}

	publicWitness, err := witness.Public()
	if err != nil {
		return gnark.Groth16Proof{}, err
	}

	// generate the proof
	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		return gnark.Groth16Proof{}, err
	}

	// verify the proof
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		return gnark.Groth16Proof{}, err
	}

	var proofBuf bytes.Buffer
	proof.WriteTo(&proofBuf)
	var vkBuf bytes.Buffer
	vk.WriteTo(&vkBuf)
	var publicWitnessBuf bytes.Buffer
	publicWitness.WriteTo(&publicWitnessBuf)
	return gnark.Groth16Proof{
		Proof:         proofBuf.Bytes(),
		VerifyingKey:  vkBuf.Bytes(),
		PublicWitness: publicWitnessBuf.Bytes(),
	}, nil
}

func TestExecuteStateChangesGroth16ECDSA(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ECDSA, takes a minute on my machine.")
	}

	f := initFixture(t)
	require := require.New(t)

	privKey, _ := ecdsa.GenerateKey(rand.Reader)
	pubkeyBytes := privKey.PublicKey.A.RawBytes()
	ethAddress := hex.EncodeToString(nativesha3.Keccak256(pubkeyBytes[:])[12:]) + ".ecdsa"

	txt := []byte("testing ECDSA (sha256)")
	sigBin, err := sign(privKey, txt)
	require.NoError(err)

	circuit, err := makeECDSACircuit(privKey, ethAddress, sigBin, txt)
	require.NoError(err)

	proof, err := generate_ecdsa_proof(circuit)
	require.NoError(err)
	jsonproof, _ := json.Marshal(proof)

	// Register the contract
	contract := zktx.Contract{
		Verifier:    "gnark-groth16-te-BN254",
		StateDigest: []byte{0},
		ProgramId:   proof.VerifyingKey,
	}

	// set the contract state
	err = f.k.Contracts.Set(f.ctx, "ecdsa", contract)
	require.NoError(err)

	msg := &zktx.MsgPublishPayloads{
		Payloads: []*zktx.Payload{
			{
				ContractName: "ecdsa",
				Data:         sigBin,
			},
		},
	}

	f.ctx = f.ctx.WithTxBytes([]byte("FakeTx"))
	h := sha256.New()
	h.Write(f.ctx.TxBytes())
	txHash := h.Sum(nil)

	_, err = f.msgServer.PublishPayloads(f.ctx, msg)
	require.NoError(err)

	msg2 := &zktx.MsgPublishPayloadProof{
		TxHash:       txHash,
		PayloadIndex: 0,
		ContractName: "ecdsa",
		Proof:        jsonproof,
	}
	_, err = f.msgServer.PublishPayloadProof(f.ctx, msg2)
	require.NoError(err)

	st, _ := f.k.Contracts.Get(f.ctx, "ecdsa")
	require.Equal(st.StateDigest, []byte{0})
}

func TestBadPayloadGroth16ECDSA(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ECDSA, takes a minute on my machine.")
	}

	require := require.New(t)

	privKey, _ := ecdsa.GenerateKey(rand.Reader)
	pubkeyBytes := privKey.PublicKey.A.RawBytes()
	ethAddress := hex.EncodeToString(nativesha3.Keccak256(pubkeyBytes[:])[12:]) + ".ecdsa"

	txt := []byte("testing ECDSA (sha256)")
	sigBin, err := sign(privKey, txt)
	require.NoError(err)

	circuit, err := makeECDSACircuit(privKey, ethAddress, sigBin, txt)
	// corrupt payload hash
	for i := 0; i < len(circuit.HyleCircuit.PayloadHash); i++ {
		circuit.HyleCircuit.PayloadHash[i] = uints.NewU8(0)
	}

	_, err = generate_ecdsa_proof(circuit)
	require.Error(err, "payload hash should mismatch")
}
