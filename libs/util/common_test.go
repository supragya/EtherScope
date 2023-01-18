package util

import (
	"encoding/hex"
	"math/big"
	"testing"
)

func TestExtractIntFromBytesPositive(t *testing.T) {
	// https://etherscan.io/tx/0xb4b7e0674e8f04e05acf1fd0fba76c512667775a11eb6d9bc83fd85f2ef47f99#eventlog lgidx 398 amount0
	_bytes, err := hex.DecodeString("0000000000000000000000000000000000000000000000000b1a2bc2ec500000")
	if err != nil {
		t.Error(err)
	}

	correctAnswer, success := big.NewInt(0).SetString("800000000000000000", 10)
	if !success {
		t.Fail()
	}

	testAnswer := ExtractIntFromBytes(_bytes)

	if testAnswer.Cmp(correctAnswer) != 0 {
		t.Error("unequal", correctAnswer, testAnswer)
	}
}

func TestExtractIntFromBytesNegative(t *testing.T) {
	// https://etherscan.io/tx/0xb4b7e0674e8f04e05acf1fd0fba76c512667775a11eb6d9bc83fd85f2ef47f99#eventlog lgidx 398 amount1
	_bytes, err := hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffb009ae26")
	if err != nil {
		t.Error(err)
	}

	correctAnswer, success := big.NewInt(0).SetString("-1341542874", 10)
	if !success {
		t.Fail()
	}

	testAnswer := ExtractIntFromBytes(_bytes)

	if testAnswer.Cmp(correctAnswer) != 0 {
		t.Error("unequal", correctAnswer, testAnswer)
	}
}

func TestExtractIntFromBytesFuzzy(t *testing.T) {
	cases := [][]string{
		{"151127302895714778", "0000000000000000000000000000000000000000000000000218e97b34f641da"},
		{"-202071", "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffceaa9"},
		{"1341542874", "000000000000000000000000000000000000000000000000000000004ff651da"},
	}

	for _, _case := range cases {
		_bytes, err := hex.DecodeString(_case[1])
		if err != nil {
			t.Error(err)
		}

		correctAnswer, success := big.NewInt(0).SetString(_case[0], 10)
		if !success {
			t.Fail()
		}

		testAnswer := ExtractIntFromBytes(_bytes)

		if testAnswer.Cmp(correctAnswer) != 0 {
			t.Error("unequal", correctAnswer, testAnswer)
		}
	}
}
