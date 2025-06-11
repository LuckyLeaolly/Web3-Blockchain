package models

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/dgraph-io/badger/v3"
)

const dbPath = "../data/blockchain"
const genesisCoinbaseData = "创世区块奖励交易"

// Blockchain 表示一个区块链
type Blockchain struct {
	lastHash []byte
	db       *badger.DB
}

// DbExists 检查数据库是否已存在
func DbExists() bool {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// InitBlockchain 创建一个新的区块链数据库
func InitBlockchain(address string) *Blockchain {
	if DbExists() {
		fmt.Println("区块链已存在")
		os.Exit(1)
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)
		fmt.Printf("创世区块创建完成: %x\n", genesis.Hash)

		err = txn.Set(genesis.Hash, genesis.Serialize())
		if err != nil {
			return err
		}

		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash

		return err
	})

	if err != nil {
		log.Panic(err)
	}

	blockchain := Blockchain{lastHash, db}

	return &blockchain
}

// NewBlockchain 返回一个现有区块链的句柄
func NewBlockchain() *Blockchain {
	if !DbExists() {
		fmt.Println("没有找到区块链，请先创建一个!")
		os.Exit(1)
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			return err
		}

		lastHash, err = item.ValueCopy(nil)
		return err
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{lastHash, db}

	return &bc
}

// MineBlock 挖掘一个新的区块
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("错误：无效的交易!")
		}
	}

	err := bc.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			return err
		}

		lastHash, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		item, err = txn.Get(lastHash)
		if err != nil {
			return err
		}

		blockData, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		block := DeserializeBlock(blockData)
		lastHeight = block.Height

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash, lastHeight+1)

	err = bc.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}

		err = txn.Set([]byte("lh"), newBlock.Hash)
		if err != nil {
			return err
		}

		bc.lastHash = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return newBlock
}

// FindUTXO 找到所有未花费的交易输出
func (bc *Blockchain) FindUTXO() map[string][]TXOutput {
	utxos := make(map[string][]TXOutput)
	spentTXOs := make(map[string][]int)

	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// 检查输出是否已经被花费
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				utxos[txID] = append(utxos[txID], out)
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return utxos
}

// FindSpendableOutputs 查找可用于交易的输出
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey([]byte(address)) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// FindUnspentTransactions 返回地址下所有未花费的交易
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)

	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// 检查输出是否已经被花费
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				if out.IsLockedWithKey([]byte(address)) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if bytes.Equal(in.PubKey, []byte(address)) { // 简化版，实际应该比较公钥哈希
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

// Iterator 返回区块链的迭代器
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.lastHash, bc.db}
}

// FindTransaction 通过ID查找交易
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("交易不存在")
}

// SignTransaction 对交易进行签名
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(vin.Txid)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction 验证交易签名
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(vin.Txid)] = prevTX
	}

	return tx.Verify(prevTXs)
}

// GetLastHash 获取区块链的最新区块哈希
func (bc *Blockchain) GetLastHash() []byte {
	return bc.lastHash
}

// GetDB 获取区块链的数据库
func (bc *Blockchain) GetDB() *badger.DB {
	return bc.db
}

// BlockchainIterator 用于遍历区块链
type BlockchainIterator struct {
	currentHash []byte
	db          *badger.DB
}

// Next 返回区块链中的下一个区块
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(i.currentHash)
		if err != nil {
			return err
		}

		blockData, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		block = DeserializeBlock(blockData)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}
