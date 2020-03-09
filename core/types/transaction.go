// Copyright 2020 The TauCoin Authors
// This file is part of the TauCoin library.
//
// The TauCoin library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The TauCoin library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
// maintained by likeopen
package types

import (
	//"container/heap"
	"errors"
	"io"
	"math/big"
	"sync/atomic"

	"github.com/Tau-Coin/taucoin-mobile-mining-go/common"
	"github.com/Tau-Coin/taucoin-mobile-mining-go/rlp"
)

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
)

type Transactiondata struct {
	tx   Transaction
	from atomic.Value
	to   atomic.Value
	hash atomic.Value
}

//define interface stands for transaction in tau
type Transaction interface {
	ChainId() Byte32s
	Protected() bool
	isProtectedV(V *big.Int) bool
	EncodeRLP(w io.Writer) error
	DecodeRLP(s *rlp.Stream) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(input []byte) error
	Fee() *big.Int
	Value() *big.Int
	Nonce() uint64
	CheckNonce() bool
	//to address
	To() *common.Address
	//get finger script
	Hash() common.Hash
	Size() common.StorageSize
	//this interface need to be repaired
	AsMessage(s Signer) (Message, error)
	WithSignature(singer Signer, sig []byte) (bool, error)
	Cost() *big.Int
	RawSignatureValues() (v, r, s *big.Int)
	GetFrom() atomic.Value
	GetSigV() *big.Int
	GetSigR() *big.Int
	GetSigS() *big.Int
	GetNounce() uint64
	GetFee() uint64
	GetReceiver() common.Address
	GetAmount() big.Int
}

//func NewTransaction(version OneByte, option OneByte, chainid Byte32s, nonce uint64, timestamp uint32, fee *big.Int, sender common.Address, receiver common.Address, amount *big.Int) *Transaction {
func NewTransaction(args ...interface{}) Transaction {
	if v, ok := args[0].(int); ok {
		//v == 0 represents transfer tx
		if v == 0 {
			return NewTransferTransaction(args[1].(OneByte),
				args[2].(OneByte),
				args[3].(Byte32s),
				args[4].(uint64),
				args[5].(uint32),
				args[6].(*big.Int),
				args[7].(common.Address),
				args[8].(common.Address),
				args[9].(*big.Int))
		}
		//v == 1 represents personal info tx
		if v == 1 {

		}
		//v == 2 represents new message tx
		if v == 2 {

		}
		//v == 3 represents new chain tx
		if v == 3 {

		}
	}
	return nil
}

type Transactions []*Transaction

func (s Transactions) Len() int { return len(s) }

func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

func TxDifference(a, b Transactions) Transactions {
	keep := make(Transactions, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, tx := range b {
		remove[(*tx).Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[(*tx).Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

//these messages need to define to adapt new ipfs system.
type Message struct {
	from       common.Address
	to         *common.Address
	nonce      uint64
	amount     *big.Int
	fee        *big.Int
	checkNonce bool
}

func NewMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, fee *big.Int, checkNonce bool) Message {
	return Message{
		from:       from,
		to:         to,
		nonce:      nonce,
		amount:     amount,
		fee:        fee,
		checkNonce: checkNonce,
	}
}

func (m Message) From() common.Address { return m.from }
func (m Message) To() *common.Address  { return m.to }
func (m Message) Nonce() uint64        { return m.nonce }
func (m Message) Value() *big.Int      { return m.amount }
func (m Message) Fee() *big.Int        { return m.fee }
func (m Message) CheckNonce() bool     { return m.checkNonce }
