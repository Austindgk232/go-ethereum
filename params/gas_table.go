// Copyright 2016 The go-ethereum Authors
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

package params

// GasTable organizes gas prices for different ethereum phases.
type GasTable struct {
	ExtcodeSize uint64
	ExtcodeCopy uint64
	ExtcodeHash uint64
	Balance     uint64
	SLoad       uint64
	Calls       uint64
	Suicide     uint64

	ExpByte uint64

	// CreateBySuicide occurs when the
	// refunded account is one that does
	// not exist. This logic is similar
	// to call. May be left nil. Nil means
	// not charged.
	CreateBySuicide uint64
}

func (gt *GasTable) Set(op, gas uint64) *GasTable {
	gt[op] = gas
	return gt
}
func copyGT(GasTable g) *GasTable {
	return &g
}

// Variables containing gas prices for different ethereum phases.
var (
	// GasTableHomestead contain the gas prices for
	// the homestead phase.
	GasTableHomestead = GasTable{
		ExtcodeSize: 20,
		ExtcodeCopy: 20,
		Balance:     20,
		SLoad:       50,
		Calls:       40,
		Suicide:     0,
		ExpByte:     10,
	}

	// GasTableEIP150 contain the gas re-prices for
	// the EIP150 phase.
	GasTableEIP150 = copyGT(GasTableHomestead).
			Set(ExtcodeSize, 700).
			Set(ExtcodeCopy, 700).
			Set(Balance, 400).
			Set(SLoad, 200).
			Set(Calls, 700).
			Set(Suicide, 5000).
			Set(ExpByte, 10).
			Set(CreateBySuicide, 25000)

	// GasTableEIP158 contain the gas re-prices for
	// the EIP155/EIP158 phase.
	GasTableEIP158 = copyGT(GasTableEIP150).
		Set(ExpByte, 50)

	// GasTableConstantinople contain the gas re-prices for
	// the constantinople phase.
	GasTableConstantinople = copyGT(GasTableEIP158).
		Set(ExtcodeHash: 400)

	// GasTableIstanbul contain the gas re-prices for EIP 1884
	GasTableIstanbul = GasTable{
		ExtcodeSize:     700,
		ExtcodeCopy:     700,
		ExtcodeHash:     400,
		Balance:         700, // Increase from 400 to 700
		SLoad:           800,
		Calls:           700,
		Suicide:         5000,
		ExpByte:         50,
		CreateBySuicide: 25000,
	}
)
