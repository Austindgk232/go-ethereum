// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package snapshot

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// Account is a slim version of a state.Account, where the root and code hash
// are replaced with a nil byte slice for empty accounts.
type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     []byte
	CodeHash []byte
}

// AccountRLP converts a state.Account content into a slim snapshot version RLP
// encoded.
func AccountRLP(nonce uint64, balance *big.Int, root common.Hash, codehash []byte) []byte {
	slim := Account{
		Nonce:   nonce,
		Balance: balance,
	}
	if root != emptyRoot {
		slim.Root = root[:]
	}
	if !bytes.Equal(codehash, emptyCode[:]) {
		slim.CodeHash = codehash
	}
	data, err := rlp.EncodeToBytes(slim)
	if err != nil {
		panic(err)
	}
	return data
}

// conversionAccount is used for converting between full and slim format. When
// doing this, we can consider 'balance' as a byte array, as it has already
// been converted from big.Int into an rlp-byteslice.
type conversionAccount struct {
	Nonce    uint64
	Balance  []byte
	Root     []byte
	CodeHash []byte
}

// SlimToFull converts data on the 'slim RLP' format into the full RLP-format
func SlimToFull(data []byte) []byte {
	acc := &conversionAccount{}
	rlp.DecodeBytes(data, acc)
	if len(acc.Root) == 0 {
		acc.Root = emptyRoot[:]
	}
	if len(acc.CodeHash) == 0 {
		acc.CodeHash = emptyCode[:]
	}
	fullData, err := rlp.EncodeToBytes(acc)
	if err != nil {
		panic(err)
	}
	return fullData
}
