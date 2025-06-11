package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"time"
)

// 难度目标位数
const targetBits = 20

// 最大随机数
const maxNonce = math.MaxInt64

// ProofOfWork 表示一个工作量证明
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork 创建一个新的工作量证明
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

// prepareData 准备用于哈希的数据
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	// 确保时间戳是有效的
	timestamp := pow.block.Timestamp
	if timestamp <= 0 {
		timestamp = time.Now().Unix() // 如果时间戳无效，使用当前时间
	}

	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run 执行工作量证明算法
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("正在挖掘新的区块 (难度: %d位)...\n", targetBits)
	startTime := time.Now()

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			elapsed := time.Since(startTime)
			fmt.Printf("\r挖矿成功! 用时: %s, 哈希: %x, 随机数: %d\n", elapsed, hash, nonce)
			break
		} else {
			nonce++
			if nonce%100000 == 0 {
				fmt.Printf("\r尝试随机数: %d", nonce)
			}
		}
	}

	if nonce == maxNonce {
		fmt.Printf("\n挖矿失败: 达到最大随机数 %d\n", maxNonce)
	}

	return nonce, hash[:]
}

// Validate 验证区块是否满足工作量证明要求
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}

// IntToHex 将int64转换为字节数组
func IntToHex(num int64) []byte {
	// 更简单可靠的实现
	buff := make([]byte, 8) // 使用固定大小的字节数组
	binary.BigEndian.PutUint64(buff, uint64(num))
	return buff
}
