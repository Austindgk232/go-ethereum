// Copyright 2020 The go-ethereum Authors
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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/trie"
)

type leaf struct {
	key   common.Hash
	value []byte
}

type trieGeneratorFn func(in chan (leaf), out chan (common.Hash))

// GenerateTrieRoot takes an account iterator and reproduces the root hash.
func GenerateTrieRoot(it AccountIterator) common.Hash {
	return generateTrieRoot(it, StdGenerate)
}

func generateTrieRoot(it AccountIterator, generatorFn trieGeneratorFn) common.Hash {
	var (
		in  = make(chan leaf)        // chan to pass leaves
		out = make(chan common.Hash) // chan to collect result
		wg  sync.WaitGroup
	)
	wg.Add(1)
	go func() {
		generatorFn(in, out)
		wg.Done()
	}()
	// Feed leaves
	start := time.Now()
	logged := time.Now()
	accounts := 0
	for it.Next() {
		slimData := it.Account()

		fullData := SlimToFull(slimData)
		l := leaf{it.Hash(), fullData}
		in <- l
		if time.Since(logged) > 8*time.Second {
			log.Info("Generating trie hash from snapshot",
				"at", l.key, "accounts", accounts, "elapsed", time.Since(start))
			logged = time.Now()
		}
		accounts++
	}
	close(in)
	result := <-out
	log.Info("Generated trie hash from snapshot", "accounts", accounts, "elapsed", time.Since(start))
	wg.Wait()
	return result
}

// StdGenerate is a very basic hexary trie builder which uses the same Trie
// as the rest of geth, with no enhancements or optimizations
func StdGenerate(in chan (leaf), out chan (common.Hash)) {
	t, _ := trie.New(common.Hash{}, trie.NewDatabase(memorydb.New()))
	for leaf := range in {
		t.TryUpdate(leaf.key[:], leaf.value)
	}
	out <- t.Hash()
}
