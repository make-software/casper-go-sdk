package types

import (
	"bytes"
	"encoding/gob"
	"math/big"
	"testing"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

func TestUInt512GobEncoding(t *testing.T) {
	var initVal = big.NewInt(1234567890123456789)
	original := clvalue.NewCLUInt512(initVal).UI512

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(&original); err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	var decoded clvalue.UInt512
	decoder := gob.NewDecoder(&buf)
	if err := decoder.Decode(&decoded); err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	if decoded.Value().Uint64() != original.Value().Uint64() {
		t.Errorf("Decoded value does not match original: got %v, want %v", decoded, original)
	}
}
