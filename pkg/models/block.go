package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

// Block 代表区块链中的一个区块
type Block struct {
	Timestamp     int64          // 区块创建时间戳
	Hash          []byte         // 当前区块的哈希值
	PrevBlockHash []byte         // 前一个区块的哈希值
	Transactions  []*Transaction // 区块中包含的交易
	Nonce         int            // 工作量证明中使用的随机数
	Height        int            // 区块高度
}

// NewBlock 创建一个新的区块
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Transactions:  transactions,
		Height:        height,
		Nonce:         0,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

// NewGenesisBlock 创建创世区块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}

// HashTransactions 计算区块中所有交易的哈希值
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// Serialize 将区块序列化为字节数组
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	if err := encoder.Encode(b); err != nil {
		panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock 从字节数组反序列化为区块
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	if err := decoder.Decode(&block); err != nil {
		panic(err)
	}

	return &block
}
