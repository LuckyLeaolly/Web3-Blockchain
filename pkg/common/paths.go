package common

import (
	"os"
	"path/filepath"
)

var (
	// 根数据目录，可以通过环境变量配置
	RootDataDir = getDefaultDataDir()
)

// 获取默认数据目录
func getDefaultDataDir() string {
	// 首先尝试从环境变量获取
	if dir := os.Getenv("BLOCKCHAIN_DATA_DIR"); dir != "" {
		return dir
	}

	// 默认使用当前目录下的data文件夹
	dir, err := filepath.Abs("data")
	if err != nil {
		// 如果获取绝对路径失败，使用相对路径
		return "data"
	}
	return dir
}

// GetBlockchainDir 获取区块链数据目录
func GetBlockchainDir() string {
	dir := filepath.Join(RootDataDir, "blockchain")
	ensureDirExists(dir)
	return dir
}

// GetWalletFilePath 获取钱包文件路径
func GetWalletFilePath() string {
	dir := filepath.Join(RootDataDir, "wallets")
	ensureDirExists(dir)
	return filepath.Join(dir, "wallet.dat")
}

// EnsureDirExists 确保目录存在
func ensureDirExists(dir string) error {
	return os.MkdirAll(dir, 0755)
}
