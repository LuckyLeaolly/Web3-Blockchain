package models

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"os"

	"personal/web3-blockchain/pkg/common"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

// 使用common包获取钱包文件路径
var walletFile = common.GetWalletFilePath()

// Wallet 表示一个钱包
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// 可序列化的钱包结构
type serializableWallet struct {
	D         []byte // 私钥D值
	X         []byte // 公钥X值
	Y         []byte // 公钥Y值
	PublicKey []byte
}

// Wallets 管理多个钱包
type Wallets struct {
	Wallets map[string]*Wallet
}

// 可序列化的钱包集合
type serializableWallets struct {
	Wallets map[string]*serializableWallet
}

// 确保钱包目录存在
func init() {
	dir := common.RootDataDir + "/wallets"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Printf("创建钱包目录失败: %v", err)
		}
	}
}

// NewWallet 创建一个新钱包
func NewWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

// NewKeyPair 创建一对新的密钥对
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

// GetAddress 从公钥获取地址
func (w *Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// 将Wallet转换为可序列化格式
func (w *Wallet) toSerializable() *serializableWallet {
	return &serializableWallet{
		D:         w.PrivateKey.D.Bytes(),
		X:         w.PrivateKey.PublicKey.X.Bytes(),
		Y:         w.PrivateKey.PublicKey.Y.Bytes(),
		PublicKey: w.PublicKey,
	}
}

// 从可序列化格式恢复Wallet
func walletFromSerializable(s *serializableWallet) *Wallet {
	curve := elliptic.P256()
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(s.D)
	priv.PublicKey.X = new(big.Int).SetBytes(s.X)
	priv.PublicKey.Y = new(big.Int).SetBytes(s.Y)

	return &Wallet{
		PrivateKey: *priv,
		PublicKey:  s.PublicKey,
	}
}

// HashPubKey 对公钥进行哈希，先用SHA-256然后用RIPEMD-160
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// Checksum 计算地址校验和
func Checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// ValidateAddress 验证地址是否有效
func ValidateAddress(address string) bool {
	pubKeyHash, err := Base58Decode([]byte(address))
	if err != nil {
		return false
	}

	if len(pubKeyHash) != 25 { // 1 byte version + 20 bytes pubKeyHash + 4 bytes checksum
		return false
	}

	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(actualChecksum, targetChecksum)
}

// GetPubKeyHashFromAddress 从地址获取公钥哈希
func GetPubKeyHashFromAddress(address string) []byte {
	pubKeyHash, err := Base58Decode([]byte(address))
	if err != nil {
		log.Panic(err)
	}

	// 去掉版本和校验和
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]

	return pubKeyHash
}

// NewWallets 从文件加载或创建钱包集合
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}

// CreateWallet 添加一个新的钱包
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

// GetAddresses 返回所有的钱包地址
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet 获取指定地址的钱包
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// LoadFromFile 从文件加载钱包
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		return err
	}

	var sWallets serializableWallets
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&sWallets)
	if err != nil {
		return err
	}

	// 转换回标准格式
	wallets := make(map[string]*Wallet)
	for addr, sWallet := range sWallets.Wallets {
		wallets[addr] = walletFromSerializable(sWallet)
	}

	ws.Wallets = wallets

	return nil
}

// Save 保存钱包到文件
func (ws *Wallets) Save() {
	// 转换为可序列化格式
	sWallets := serializableWallets{
		Wallets: make(map[string]*serializableWallet),
	}

	for addr, wallet := range ws.Wallets {
		sWallets.Wallets[addr] = wallet.toSerializable()
	}

	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(sWallets)
	if err != nil {
		log.Printf("编码钱包数据时出错: %v", err)
		log.Panic(err)
	}

	// 创建临时文件
	tmpFile := walletFile + ".tmp"
	err = os.WriteFile(tmpFile, content.Bytes(), 0644)
	if err != nil {
		log.Printf("写入钱包临时文件时出错: %v", err)
		log.Panic(err)
	}

	// 安全地重命名文件，保证原子性操作
	err = os.Rename(tmpFile, walletFile)
	if err != nil {
		log.Printf("重命名钱包文件时出错: %v", err)
		log.Panic(err)
	}

	log.Printf("钱包数据已成功保存到 %s", walletFile)
}

// Base58Encode 进行Base58编码
func Base58Encode(input []byte) []byte {
	var result []byte

	// 实际应用中应该使用完整的Base58编码算法
	// 这里为了简化，我们只实现一个基本版本
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := &big.Int{}

	// Base58字母表
	alphabet := []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = append(result, alphabet[mod.Int64()])
	}

	// 反转结果
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	// 处理前导0（在Base58中表示为1）
	for _, b := range input {
		if b == 0x00 {
			result = append([]byte{alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result
}

// Base58Decode 进行Base58解码
func Base58Decode(input []byte) ([]byte, error) {
	// Base58字母表
	alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	// 创建一个长度为256的查找表，用于快速找到各个字符在alphabet中的位置
	var lookupTable [256]byte
	for i := 0; i < len(lookupTable); i++ {
		lookupTable[i] = 255 // 设置一个无效值
	}
	for i := 0; i < len(alphabet); i++ {
		lookupTable[alphabet[i]] = byte(i)
	}

	// 计算出最终结果的最大可能长度
	var zerosCount int
	for i := 0; i < len(input); i++ {
		if input[i] != alphabet[0] {
			break
		}
		zerosCount++
	}

	// 准备解码后数据的存储空间
	result := big.NewInt(0)
	base := big.NewInt(58)

	// 从左到右逐个处理字符
	for i := zerosCount; i < len(input); i++ {
		charIndex := lookupTable[input[i]]
		if charIndex == 255 {
			return nil, fmt.Errorf("非法字符 '%c' 在位置 %d", input[i], i)
		}

		result.Mul(result, base)
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	// 将大整数转换为字节数组
	decoded := result.Bytes()

	// 处理前导零
	decodedLength := zerosCount + len(decoded)
	resultBytes := make([]byte, decodedLength)
	copy(resultBytes[zerosCount:], decoded)

	return resultBytes, nil
}
