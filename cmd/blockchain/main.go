package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"personal/web3-blockchain/pkg/models"
)

// CLI 交互式命令行
type CLI struct{}

// printUsage 显示命令行用法
func (cli *CLI) printUsage() {
	fmt.Println("用法:")
	fmt.Println("  createblockchain -address ADDRESS - 创建一个新的区块链并发送创世区块奖励到ADDRESS")
	fmt.Println("  createwallet - 创建一个新的钱包")
	fmt.Println("  getbalance -address ADDRESS - 获取地址余额")
	fmt.Println("  listaddresses - 列出所有钱包地址")
	fmt.Println("  printchain - 打印区块链中的所有区块")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - 发送AMOUNT的币从FROM地址到TO地址")
}

// validateArgs 验证命令行参数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

// createBlockchain 创建区块链
func (cli *CLI) createBlockchain(address string) {
	if !models.ValidateAddress(address) {
		log.Panic("错误: 地址非法")
	}
	bc := models.InitBlockchain(address)
	defer bc.GetDB().Close()

	fmt.Println("完成!")
}

// getBalance 获取地址余额
func (cli *CLI) getBalance(address string) {
	if !models.ValidateAddress(address) {
		log.Panic("错误: 地址非法")
	}

	bc := models.NewBlockchain()
	defer bc.GetDB().Close()

	balance := 0
	pubKeyHash := models.GetPubKeyHashFromAddress(address)
	UTXOs := bc.FindUTXO()

	for _, outputs := range UTXOs {
		for _, output := range outputs {
			if output.IsLockedWithKey(pubKeyHash) {
				balance += output.Value
			}
		}
	}

	fmt.Printf("'%s'的余额: %d\n", address, balance)
}

// printChain 打印区块链
func (cli *CLI) printChain() {
	bc := models.NewBlockchain()
	defer bc.GetDB().Close()

	iter := bc.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("前一个区块哈希: %x\n", block.PrevBlockHash)
		fmt.Printf("当前区块哈希: %x\n", block.Hash)
		pow := models.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", pow.Validate())

		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

// send 发送交易
func (cli *CLI) send(from, to string, amount int) {
	if !models.ValidateAddress(from) {
		log.Panic("错误: 发送地址非法")
	}
	if !models.ValidateAddress(to) {
		log.Panic("错误: 接收地址非法")
	}

	bc := models.NewBlockchain()
	defer bc.GetDB().Close()

	tx := models.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*models.Transaction{tx})
	fmt.Println("成功发送!")
}

// Run 运行CLI
func (cli *CLI) Run() {
	cli.validateArgs()

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "区块链创建者地址")
	getBalanceAddress := getBalanceCmd.String("address", "", "要查询余额的地址")
	sendFrom := sendCmd.String("from", "", "源钱包地址")
	sendTo := sendCmd.String("to", "", "目标钱包地址")
	sendAmount := sendCmd.Int("amount", 0, "发送数量")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

// createWallet 创建钱包
func (cli *CLI) createWallet() {
	wallets, _ := models.NewWallets()
	address := wallets.CreateWallet()
	wallets.Save()

	fmt.Printf("你的新地址: %s\n", address)
}

// listAddresses 列出所有钱包地址
func (cli *CLI) listAddresses() {
	wallets, err := models.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func main() {
	defer os.Exit(0)
	cli := CLI{}
	cli.Run()
}
