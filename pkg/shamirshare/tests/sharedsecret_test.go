package tests

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/shamir"
	"gopkg.in/go-playground/assert.v1"
)

var encoder *base64.Encoding = base64.URLEncoding

func TestDistribute(t *testing.T) {
	secret := []byte("This is my secret!!")
	nb_shards := 3
	min_nb_shards_needed := 2
	keys, err := shamir.Split(secret, nb_shards, min_nb_shards_needed)
	if err != nil {
		t.Errorf("error to udpdate keys: %v", err)
	}
	assert.Equal(t, len(keys) == nb_shards, true)
	for _, shard := range keys {
		encoded := encoder.EncodeToString(shard)
		fmt.Println(encoded)
	}
}

func TestCombine(t *testing.T) {
	secret := "This is my secret!!"
	nb_shards := 3
	min_nb_shards_needed := 2
	shardBytes, _ := shamir.Split([]byte(secret), nb_shards, min_nb_shards_needed)
	decodedSecret, err := shamir.Combine(shardBytes)
	if err != nil {
		t.Errorf("error to udpdate keys: %v", err)
	}
	assert.Equal(t, string(decodedSecret) == secret, true)

	shards := [][]byte{shardBytes[0], shardBytes[1]}
	decodedSecret, err = shamir.Combine(shards)
	if err != nil {
		t.Errorf("error to udpdate keys: %v", err)
	}
	assert.Equal(t, string(decodedSecret) == secret, true)

	shards = [][]byte{shardBytes[1], shardBytes[2]}
	decodedSecret, err = shamir.Combine(shards)
	if err != nil {
		t.Errorf("error to udpdate keys: %v", err)
	}
	assert.Equal(t, string(decodedSecret) == secret, true)

	shards = [][]byte{shardBytes[0], shardBytes[2]}
	decodedSecret, err = shamir.Combine(shards)
	if err != nil {
		t.Errorf("error to udpdate keys: %v", err)
	}
	assert.Equal(t, string(decodedSecret) == secret, true)

}
