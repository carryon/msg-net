package http_rpc_server

import (
	"math/big"
	"sync/atomic"
	"encoding/hex"
	"fmt"
	"github.com/bocheninc/msg-net/utils"
)

const (
	addressPrefix = "0x"
	AddressLength = 20
	// HashSize represents the hash length
	HashSize = 32
	// SignatureSize represents the signature length
	SignatureSize = 65
)

type Address [AddressLength]byte

type message struct {
	ChainId   string  `json:"chainId"`
}

// Message represents the message transfer in msg-net
type sendMessage struct {
	Cmd     uint8
	Payload []byte
}

type transformMiddleData struct {
	CmdType uint8
	Payload []byte
}

type test_struct struct {
	Test string
}

type (
	// Hash represents the 32 byte hash of arbitrary data
	Hash [HashSize]byte
	Signature [SignatureSize]byte
	ChainCoordinate []byte
)

type pRouterHandler interface {
	IsPeerExist(id string) bool
}


type BroadcastReply struct {
	ContractAddr    *string     `json:"contractAddr"`
	TransactionHash Hash       `json:"transactionHash"`
}

type Balance struct {
	Amount *big.Int
	Nonce  uint32
}

// BlockHeader represents the header in block
type BlockHeader struct {
	PreviousHash  Hash        `json:"previousHash" `
	TimeStamp     uint32      `json:"timeStamp"`
	Nonce         uint32      `json:"nonce" `
	TxsMerkleHash Hash        `json:"transactionsMerkleHash" `
	Height        uint32      `json:"height" `
}

//Block json rpc return block
type Block struct {
	BlockHeader *BlockHeader `json:"header"`
	TxHashList  []Hash      `json:"txHashList"`
}

type txdata struct {
	FromChain  ChainCoordinate `json:"fromChain"`
	ToChain    ChainCoordinate `json:"toChain"`
	Type       uint32                     `json:"type"`
	Nonce      uint32                     `json:"nonce"`
	Sender     Address           `json:"sender"`
	Recipient  Address           `json:"recipient"`
	Amount     *big.Int                   `json:"amount"`
	Fee        *big.Int                   `json:"fee"`
	Signature  *Signature          `json:"signature"`
	CreateTime uint32                     `json:"createTime"`
}

// Transaction represents the basic transaction that contained in blocks
type Transaction struct {
	Data    txdata `json:"data"`
	Payload []byte `json:"payload"`

	hash   atomic.Value
	sender atomic.Value
}

// Transactions represents transaction slice type for basic sorting.
type Transactions []*Transaction

type SerialDataType []byte

// String returns address string
func (self Address) String() string {
	return fmt.Sprintf("%s%x", addressPrefix, self[:])
}


//Bytes returns address bytes
func (self Address) Bytes() []byte {
	return self[:]
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalText() ([]byte, error) {
	return utils.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *Address) UnmarshalText(input []byte) error {
	return utils.UnmarshalFixedText(input, a[:])
}

func (h Hash) String() string { return hex.EncodeToString(h[:]) }

// MarshalText returns the hex representation of h.
func (h Hash) MarshalText() ([]byte, error) {
	return utils.Bytes(h[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	return utils.UnmarshalFixedText(input, h[:])
}
