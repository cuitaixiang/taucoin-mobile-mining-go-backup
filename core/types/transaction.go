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
	"container/heap"
	"io"
	"errors"
	"math/big"
	"sync/atomic"

	"github.com/Tau-Coin/taucoin-mobile-mining-go/common"
	"github.com/Tau-Coin/taucoin-mobile-mining-go/common/hexutil"
	"github.com/Tau-Coin/taucoin-mobile-mining-go/crypto"
	"github.com/Tau-Coin/taucoin-mobile-mining-go/rlp"
	"golang.org/x/crypto/sha3"
)

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
)

//go:generate gencodec -type txdata -field-override txdataMarshaling -out transfer_tx_json.go
type Transaction struct {
	data txdata

	hash atomic.Value
	size atomic.Value
	from atomic.Value
}
type txdata struct {
	Version   OneByte         `json:"version"     gencodec:"required"`
	Option    OneByte         `json:"option"      gencodec:"required"`
	ChainID   Byte32s         `json:"chainid"     gencodec:"required"`
	Nonce    uint64          `json:"nonce"      gencodec:"required"`
	TimeStamp uint32          `json:"timestamp"   gencodec:"required"`
	//Fee       OneByte         `json:"fee"         gencodec:"required"`
	Fee       *big.Int        `json:"fee"         gencodec:"required"`
	V         *big.Int        `json:"v"           gencodec:"required"`
	R         *big.Int        `json:"r"           gencodec:"required"`
	S         *big.Int        `json:"s"           gencodec:"required"`
	Sender    *common.Address `json:"sender"        rlp:"required"`

	Receiver *common.Address `json:"receiver"        rlp:"required"`
	//Amount   Byte5s          `json:"amount"       gencodec:"required"`
	Amount     *big.Int        `json:"value"    gencodec:"required"`
}

type txdataMarshaling struct {
	Version   hexutil.Bytes
	Option    hexutil.Bytes
	ChainID   hexutil.Bytes
	Nonce    hexutil.Uint64
	TimeStamp hexutil.Uint32
	//Fee       hexutil.Bytes
	Fee       *hexutil.Big
	V         *hexutil.Big
	R         *hexutil.Big
	S         *hexutil.Big

	//Amount hexutil.Bytes
	Amount       *hexutil.Big
}

//func NewTransaction(version OneByte, option OneByte, chainid Byte32s, nonce uint64, timestamp uint32, fee *big.Int, sender common.Address, receiver common.Address, amount *big.Int) *Transaction {
func NewTransaction(nonce uint64, receiver common.Address, amount *big.Int, fee *big.Int) *Transaction {
	var version OneByte
	var option OneByte
	var chainid Byte32s
	var timestamp uint32
	var sender common.Address
	return newTransaction(version, option, chainid, nonce, timestamp, fee, &sender, &receiver, amount)
}

func newTransaction(version OneByte, option OneByte, chainid Byte32s, nonce uint64, timestamp uint32, fee *big.Int, sender *common.Address, receiver *common.Address, amount *big.Int) *Transaction {
	d := txdata{
		Version:   version,
		Option:    option,
		ChainID:   chainid,
		Nonce:    nonce,
		TimeStamp: timestamp,
		Fee:       fee,
		V:         new(big.Int),
		R:         new(big.Int),
		S:         new(big.Int),
		Sender:    sender,
		Receiver:  receiver,
		//Amount:    amount,
		Amount:    new(big.Int),
	}
	if amount != nil {
		d.Amount.Set(amount)
	}
	return &Transaction{data: d}
}

func (tx *Transaction) ChainId() Byte32s {
	return tx.data.ChainID
}

func (tx *Transaction) Protected() bool {
	return isProtectedV(tx.data.V)
}

func isProtectedV(V *big.Int) bool {
	v := V.Uint64()
	if v == 27 || v == 28 {
		return false
	}

	return true
}

func (tx *Transaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &tx.data)
}

func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&tx.data)
	if err == nil {
		tx.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	data := tx.data
	return data.MarshalJSON()
}

func (tx *Transaction) UnmarshalJSON(input []byte) error {
	var dec txdata
	if err := dec.UnmarshalJSON(input); err != nil {
		return err
	}

	withSignature := dec.V.Sign() != 0 || dec.R.Sign() != 0 || dec.S.Sign() != 0
	if withSignature {
		var V byte
		if isProtectedV(dec.V) {
			chainID := deriveChainId(dec.V).Uint64()
			V = byte(dec.V.Uint64() - 35 - 2*chainID)
		} else {
			V = byte(dec.V.Uint64() - 27)
		}
		if !crypto.ValidateSignatureValues(V, dec.R, dec.S, false) {
			return ErrInvalidSig
		}
	}

	*tx = Transaction{data: dec}
	return nil
}

func (tx *Transaction) Fee() *big.Int    { return new(big.Int).Set(tx.data.Fee)}
//func (tx *Transaction) Value() Byte5s  { return tx.data.Amount }
func (tx *Transaction) Value() *big.Int  { return new(big.Int).Set(tx.data.Amount)}
func (tx *Transaction) Nonce() uint64    { return tx.data.Nonce }
func (tx *Transaction) CheckNonce() bool { return true }

// To returns the recipient address of the transaction.
// It returns nil if the transaction is a contract creation.
func (tx *Transaction) To() *common.Address {
	if tx.data.Receiver == nil {
		return nil
	}
	to := *tx.data.Receiver
	return &to
}


func (tx *Transaction) Hash() (h common.Hash) {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, tx)
	hw.Sum(h[:0])

	tx.hash.Store(h)
	return h
}

func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}

	c := writeCounter(0)
	rlp.Encode(&c, &tx.data)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

func (tx *Transaction) AsMessage(s Signer) (Message, error) {
	msg := Message{
		//todo
		nonce:      tx.data.Nonce,
		fee:        tx.data.Fee,
		to:         tx.data.Receiver,
		amount:     tx.data.Amount,
		checkNonce: true,
	}

	var err error
	//msg.from, err = Sender(s, data)
	return msg, err
}

func (tx *Transaction) WithSignature(singer Signer, sig []byte) (*Transaction, error) {
	//todo splite and verify data
	cpy := &Transaction{data: tx.data}
	//fill field of signature in data
	return cpy, nil
}

func (tx *Transaction) Cost() *big.Int {
	cost := tx.data.Amount
	fee := tx.data.Fee
	return cost.Add(cost, fee)
}

func (tx *Transaction) RawSignatureValues() (v, r, s *big.Int) {
	return tx.data.V, tx.data.R, tx.data.S
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
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

type TxByNonce Transactions

func (s TxByNonce) Len() int { return len(s) }
func (s TxByNonce) Less(i, j int) bool {
	return s[i].data.Nonce < s[j].data.Nonce
}
func (s TxByNonce) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s *TxByNonce) Push(x interface{}) {
	*s = append(*s, x.(*Transaction))
}

func (s *TxByNonce) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

type TxByPrice Transactions

func (s TxByPrice) Len() int { return len(s) }
func (s TxByPrice) Less(i, j int) bool {
	fee1 := new(big.Int).Set(s[i].data.Fee)
	fee2 := new(big.Int).Set(s[j].data.Fee)
	return fee1.Cmp(fee2) > 0
}
func (s TxByPrice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s *TxByPrice) Push(x interface{}) {
	*s = append(*s, x.(*Transaction))
}

func (s *TxByPrice) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

type TransactionsByFeeAndNonce struct {
	txs    map[common.Address]Transactions
	heads  TxByPrice
	signer Signer
}

//watch out Transactions is sorted by Nonce first
func NewTransactionsByFeeAndNonce(signer Signer, txs map[common.Address]Transactions) *TransactionsByFeeAndNonce {
	heads := make(TxByPrice, 0, len(txs))
	for _, accTxs := range txs {
		heads = append(heads, accTxs[0])
		//to make sure a list txs from txs is from same account
		//to do to complete singer
	}
	heap.Init(&heads)
	return &TransactionsByFeeAndNonce{
		txs:   txs,
		heads: heads,
		//This singer need to make adaption mpdify abount unique tx
		signer: signer,
	}
}

func (t *TransactionsByFeeAndNonce) Peek() *Transaction {
	if len(t.heads) == 0 {
		return nil
	}

	return t.heads[0]
}

func (t *TransactionsByFeeAndNonce) Shift() {
	//if Account x contains other txs sorted by Nonce, the others should
	//come up with new fee sorting because of some fee element changing.
	acc := t.heads[0].data.Sender
	if txs, ok := t.txs[*acc]; ok && len(txs) > 0 {
		t.heads[0], t.txs[*acc] = txs[0], txs[1:]
		heap.Fix(&t.heads, 0)
	} else {
		heap.Pop(&t.heads)
	}
}

func (t *TransactionsByFeeAndNonce) Pop() {
	heap.Pop(&t.heads)
}

//todo
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
