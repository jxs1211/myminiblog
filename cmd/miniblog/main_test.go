package main

import (
	"testing"

	"github.com/bytedance/sonic"
)

func TestSonic(t *testing.T) {
	var data struct{ Name string }
	// Marshal
	output, err := sonic.Marshal(&data)
	if err != nil {
		t.Fatal(err)
	}
	// Unmarshal
	err = sonic.Unmarshal(output, &data)
	if err != nil {
		t.Fatal(err)
	}
}
