// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/Tau-Coin/taucoin-mobile-mining-go/common"
	"github.com/Tau-Coin/taucoin-mobile-mining-go/common/hexutil"
)

var _ = (*PersonalInfoTxDataMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (p PersonalInfoTxData) MarshalJSON() ([]byte, error) {
	type PersonalInfoTxData struct {
		Version     hexutil.Bytes   `json:"version"     gencodec:"required"`
		Option      hexutil.Bytes   `json:"option"      gencodec:"required"`
		ChainID     hexutil.Bytes   `json:"chainid"     gencodec:"required"`
		Nonce       hexutil.Uint64  `json:"nonce"      gencodec:"required"`
		TimeStamp   hexutil.Uint32  `json:"timestamp"   gencodec:"required"`
		Fee         *hexutil.Big    `json:"fee"         gencodec:"required"`
		V           *hexutil.Big    `json:"v"           gencodec:"required"`
		R           *hexutil.Big    `json:"r"           gencodec:"required"`
		S           *hexutil.Big    `json:"s"           gencodec:"required"`
		Sender      *common.Address `json:"sender"        rlp:"required"`
		ContactName hexutil.Bytes   `json:"contactname" gencodec:"required"`
		Name        hexutil.Bytes   `json:"name"        gencodec:"required"`
		Profile     hexutil.Bytes   `json:"profile"     gencodec:"required"`
	}
	var enc PersonalInfoTxData
	enc.Version = hexutil.Bytes(p.Version)
	enc.Option = hexutil.Bytes(p.Option)
	enc.ChainID = hexutil.Bytes(p.ChainID)
	enc.Nonce = hexutil.Uint64(p.Nonce)
	enc.TimeStamp = hexutil.Uint32(p.TimeStamp)
	enc.Fee = (*hexutil.Big)(p.Fee)
	enc.V = (*hexutil.Big)(p.V)
	enc.R = (*hexutil.Big)(p.R)
	enc.S = (*hexutil.Big)(p.S)
	enc.Sender = p.Sender
	enc.ContactName = hexutil.Bytes(p.ContactName)
	enc.Name = hexutil.Bytes(p.Name)
	enc.Profile = hexutil.Bytes(p.Profile)
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (p *PersonalInfoTxData) UnmarshalJSON(input []byte) error {
	type PersonalInfoTxData struct {
		Version     *hexutil.Bytes  `json:"version"     gencodec:"required"`
		Option      *hexutil.Bytes  `json:"option"      gencodec:"required"`
		ChainID     *hexutil.Bytes  `json:"chainid"     gencodec:"required"`
		Nonce       *hexutil.Uint64 `json:"nonce"      gencodec:"required"`
		TimeStamp   *hexutil.Uint32 `json:"timestamp"   gencodec:"required"`
		Fee         *hexutil.Big    `json:"fee"         gencodec:"required"`
		V           *hexutil.Big    `json:"v"           gencodec:"required"`
		R           *hexutil.Big    `json:"r"           gencodec:"required"`
		S           *hexutil.Big    `json:"s"           gencodec:"required"`
		Sender      *common.Address `json:"sender"        rlp:"required"`
		ContactName *hexutil.Bytes  `json:"contactname" gencodec:"required"`
		Name        *hexutil.Bytes  `json:"name"        gencodec:"required"`
		Profile     *hexutil.Bytes  `json:"profile"     gencodec:"required"`
	}
	var dec PersonalInfoTxData
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Version == nil {
		return errors.New("missing required field 'version' for PersonalInfoTxData")
	}
	p.Version = OneByte(*dec.Version)
	if dec.Option == nil {
		return errors.New("missing required field 'option' for PersonalInfoTxData")
	}
	p.Option = OneByte(*dec.Option)
	if dec.ChainID == nil {
		return errors.New("missing required field 'chainid' for PersonalInfoTxData")
	}
	p.ChainID = Byte32s(*dec.ChainID)
	if dec.Nonce == nil {
		return errors.New("missing required field 'nonce' for PersonalInfoTxData")
	}
	p.Nonce = uint64(*dec.Nonce)
	if dec.TimeStamp == nil {
		return errors.New("missing required field 'timestamp' for PersonalInfoTxData")
	}
	p.TimeStamp = uint32(*dec.TimeStamp)
	if dec.Fee == nil {
		return errors.New("missing required field 'fee' for PersonalInfoTxData")
	}
	p.Fee = (*big.Int)(dec.Fee)
	if dec.V == nil {
		return errors.New("missing required field 'v' for PersonalInfoTxData")
	}
	p.V = (*big.Int)(dec.V)
	if dec.R == nil {
		return errors.New("missing required field 'r' for PersonalInfoTxData")
	}
	p.R = (*big.Int)(dec.R)
	if dec.S == nil {
		return errors.New("missing required field 's' for PersonalInfoTxData")
	}
	p.S = (*big.Int)(dec.S)
	if dec.Sender != nil {
		p.Sender = dec.Sender
	}
	if dec.ContactName == nil {
		return errors.New("missing required field 'contactname' for PersonalInfoTxData")
	}
	p.ContactName = Byte32s(*dec.ContactName)
	if dec.Name == nil {
		return errors.New("missing required field 'name' for PersonalInfoTxData")
	}
	p.Name = Byte20s(*dec.Name)
	if dec.Profile == nil {
		return errors.New("missing required field 'profile' for PersonalInfoTxData")
	}
	p.Profile = Byte32s(*dec.Profile)
	return nil
}
