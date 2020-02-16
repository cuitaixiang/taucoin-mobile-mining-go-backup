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
	"io"
	"math/big"
	"sync/atomic"

	"github.com/Tau-Coin/taucoin-mobile-mining-go/common"
	"github.com/Tau-Coin/taucoin-mobile-mining-go/common/hexutil"
	"github.com/Tau-Coin/taucoin-mobile-mining-go/crypto"
	"github.com/Tau-Coin/taucoin-mobile-mining-go/rlp"
	"golang.org/x/crypto/sha3"
)

//go:generate gencodec -type TransferTxData -field-override TransferTxDataMarshaling -out transfer_tx_json.go
type TransferTx struct {
	tx TransferTxData

	hash atomic.Value
	size atomic.Value
	from atomic.Value
}
type Byte5s []byte

type TransferTxData struct {
	Version   OneByte         `json:"version"     gencodec:"required"`
	Option    OneByte         `json:"option"      gencodec:"required"`
	ChainID   Byte32s         `json:"chainid"     gencodec:"required"`
	Nounce    uint64          `json:"nounce"      gencodec:"required"`
	TimeStamp uint32          `json:"timestamp"   gencodec:"required"`
	Fee       OneByte         `json:"fee"         gencodec:"required"`
	V         *big.Int        `json:"v"           gencodec:"required"`
	R         *big.Int        `json:"r"           gencodec:"required"`
	S         *big.Int        `json:"s"           gencodec:"required"`
	Sender    *common.Address `json:"sender"        rlp:"required"`

	Receiver *common.Address `json:"sender"        rlp:"required"`
	Amount   Byte5s          `json:"amount"       gencodec:"required"`
}

type TransferTxDataMarshaling struct {
	Version   hexutil.Bytes
	Option    hexutil.Bytes
	ChainID   hexutil.Bytes
	Nounce    hexutil.Uint64
	TimeStamp hexutil.Uint32
	Fee       hexutil.Bytes
	V         *hexutil.Big
	R         *hexutil.Big
	S         *hexutil.Big

	Amount hexutil.Bytes
}
/*
func NewTransferTransaction(version OneByte, option OneByte, chainid OneByte, nounce uint64, timestamp uint32, fee OneByte, sender common.Address, receiver common.Address, amount Byte5s) {
	return newTransferTransaction(version, option, chainid, nounce, timestamp, fee, &sender, &receiver, amount)
}
func newTransferTransaction(version OneByte, option OneByte, chainid OneByte, nounce uint64, timestamp uint32, fee OneByte, sender *common.Address, receiver *common.Address, amount Byte5s) *TransferTx {
	d := TransferTxData{
		Version:   version,
		Option:    option,
		ChainID:   chainid,
		Nounce:    nounce,
		TimeStamp: timestamp,
		Fee:       fee,
		V:         new(big.Int),
		R:         new(big.Int),
		S:         new(big.Int),
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
	}
	return TransferTx{tx: d}
}
*/

func (ttx *TransferTx) ChainId() Byte32s {
	return ttx.tx.ChainID
}

func (ttx *TransferTx) Protected() bool {
	return true
}

func (ttx *TransferTx) isProtectedV(V *big.Int) bool {
	v := V.Uint64()
	if v == 27 || v == 28 {
		return false
	}

	return true
}

func (ttx *TransferTx) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &ttx.tx)
}

func (ttx *TransferTx) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&ttx.tx)
	if err == nil {
		ttx.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

func (ttx *TransferTx) MarshalJSON() ([]byte, error) {
	data := ttx.tx
	return data.MarshalJSON()
}

func (ttx *TransferTx) UnmarshalJSON(input []byte) error {
	var dec TransferTxData
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

	*ttx = TransferTx{tx: dec}
	return nil
}

func (ttx *TransferTx) Fee() OneByte     { return ttx.tx.Fee }
func (ttx *TransferTx) Value() Byte5s    { return ttx.tx.Amount }
func (ttx *TransferTx) Nonce() uint64    { return ttx.tx.Nounce }
func (ttx *TransferTx) CheckNonce() bool { return true }

func (ttx *TransferTx) Hash() (h common.Hash) {
	if hash := ttx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, ttx)
	hw.Sum(h[:0])

	ttx.hash.Store(h)
	return h
}

func (ttx *TransferTx) Size() common.StorageSize {
	if size := ttx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}

	c := writeCounter(0)
	rlp.Encode(&c, &ttx.tx)
	ttx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

func (ttx *TransferTx) AsMessage(s Signer) (Message, error) {
	msg := Message{
		//todo
	}

	var err error
	// todo should compare derive sender with contained sender
	return msg, err
}

func (ttx *TransferTx) WithSignature(singer Signer, sig []byte) (*TransferTx, error) {
	//todo splite and verify ttx
	cpy := &TransferTx{tx: ttx.tx}
	//fill field of signature in ttx
	return cpy, nil
}

func (ttx *TransferTx) Cost() *big.Int {
	cost := new(big.Int)
	cost.SetBytes(ttx.tx.Amount)
	fee := new(big.Int)
	fee.SetBytes(ttx.tx.Fee)
	return cost.Add(cost, fee)
}

func (ttx *TransferTx) RawSignatureValues() (v, r, s *big.Int) {
	return ttx.tx.V, ttx.tx.R, ttx.tx.S
}
