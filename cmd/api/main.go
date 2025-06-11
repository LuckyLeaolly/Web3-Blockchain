package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"personal/web3-blockchain/docs"
	"personal/web3-blockchain/pkg/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// BlockchainServer 区块链API服务
type BlockchainServer struct {
	blockchain *models.Blockchain
}

// BlockResponse 区块响应结构
type BlockResponse struct {
	Hash          string              `json:"hash"`
	PrevBlockHash string              `json:"prevBlockHash"`
	Timestamp     int64               `json:"timestamp"`
	Height        int                 `json:"height"`
	Nonce         int                 `json:"nonce"`
	Transactions  []TransactionOutput `json:"transactions"`
}

// TransactionOutput 交易输出结构
type TransactionOutput struct {
	ID        string   `json:"id"`
	From      string   `json:"from"`
	To        string   `json:"to"`
	Amount    int      `json:"amount"`
	Timestamp int64    `json:"timestamp"`
	Inputs    []string `json:"inputs"`
}

// NewBlockchainServer 创建新的区块链服务
func NewBlockchainServer() *BlockchainServer {
	bc := models.NewBlockchain()
	return &BlockchainServer{blockchain: bc}
}

// @title Web3.0 区块链系统 API
// @version 1.0
// @description 这是一个使用Golang实现的Web3.0区块链系统的API服务
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.web3blockchain.io/support
// @contact.email support@web3blockchain.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

func main() {
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 创建区块链服务
	server := NewBlockchainServer()

	// API路由组
	api := r.Group("/api/v1")
	{
		// 认证相关API - 不需要认证
		auth := api.Group("/auth")
		{
			auth.POST("/login", server.login)
			auth.POST("/register", server.register)
		}

		// 公开API - 不需要认证
		api.GET("/info", server.getBlockchainInfo)
		api.GET("/blocks", server.getBlocks)
		api.GET("/blocks/:hash", server.getBlock)
		api.GET("/blocks/height/:height", server.getBlockByHeight)
		api.GET("/transactions", server.getTransactions)
		api.GET("/transactions/:id", server.getTransaction)
		api.GET("/wallets/:address/balance", server.getBalance)
		api.GET("/wallets/:address/transactions", server.getWalletTransactions)

		// 需要认证的API
		protected := api.Group("")
		protected.Use(JWTAuthMiddleware())
		{
			protected.POST("/transactions", server.createTransaction)
			protected.GET("/wallets", server.getWallets)
			protected.POST("/wallets", server.createWallet)
			protected.GET("/network/config", server.getNetworkConfig)
			protected.PUT("/network/config", server.updateNetworkConfig)
		}
	}

	// 自定义Swagger路由
	r.GET("/swagger/*any", func(c *gin.Context) {
		// 从我们的docs包获取swagger.json
		if c.Param("any") == "/swagger.json" {
			jsonData, err := docs.GetSwaggerJSON()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "无法加载Swagger文档"})
				return
			}
			c.Data(http.StatusOK, "application/json", jsonData)
			return
		}
		// 其他路径使用ginSwagger处理
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	// 启动服务
	log.Println("API服务已启动，监听端口8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}

// @Summary 获取区块链信息
// @Description 返回区块链的基本信息，如区块高度、交易总数等
// @Tags 区块链信息
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /info [get]
func (s *BlockchainServer) getBlockchainInfo(c *gin.Context) {
	// 获取区块链状态
	height := 0
	txCount := 0

	// 获取最后一个区块的高度和总交易数
	iter := s.blockchain.Iterator()
	for {
		block := iter.Next()
		height = block.Height
		txCount += len(block.Transactions)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"height":       height,
		"transactions": txCount,
		"status":       "running",
		"version":      "1.0.0",
	})
}

// @Summary 获取所有区块
// @Description 返回区块链中的所有区块
// @Tags 区块
// @Accept json
// @Produce json
// @Param limit query int false "限制返回的区块数量"
// @Param offset query int false "跳过前面的区块数量"
// @Success 200 {array} BlockResponse
// @Router /blocks [get]
func (s *BlockchainServer) getBlocks(c *gin.Context) {
	var blocks []BlockResponse

	// 解析分页参数
	limit := 10 // 默认返回10个区块
	offset := 0

	iter := s.blockchain.Iterator()
	currentBlock := 0

	for {
		block := iter.Next()

		// 应用offset跳过前面的区块
		if currentBlock >= offset {
			blocks = append(blocks, convertBlockToResponse(block))
		}

		currentBlock++

		// 如果达到limit限制或到达创世区块，则退出循环
		if (limit > 0 && len(blocks) >= limit) || len(block.PrevBlockHash) == 0 {
			break
		}
	}

	c.JSON(http.StatusOK, blocks)
}

// @Summary 获取特定区块
// @Description 通过区块哈希获取区块详情
// @Tags 区块
// @Accept json
// @Produce json
// @Param hash path string true "区块哈希"
// @Success 200 {object} BlockResponse
// @Failure 404 {object} map[string]string
// @Router /blocks/{hash} [get]
func (s *BlockchainServer) getBlock(c *gin.Context) {
	hash := c.Param("hash")

	// 查找区块
	iter := s.blockchain.Iterator()
	found := false
	var foundBlock *models.Block

	for {
		block := iter.Next()

		// 如果找到目标区块
		if fmt.Sprintf("%x", block.Hash) == hash {
			found = true
			foundBlock = block
			break
		}

		// 如果到达创世区块仍未找到，退出循环
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "区块未找到"})
		return
	}

	c.JSON(http.StatusOK, convertBlockToResponse(foundBlock))
}

// @Summary 通过高度获取区块
// @Description 通过区块高度获取区块详情
// @Tags 区块
// @Accept json
// @Produce json
// @Param height path int true "区块高度"
// @Success 200 {object} BlockResponse
// @Failure 404 {object} map[string]string
// @Router /blocks/height/{height} [get]
func (s *BlockchainServer) getBlockByHeight(c *gin.Context) {
	heightStr := c.Param("height")
	var height int

	// 解析高度参数
	if _, err := fmt.Sscanf(heightStr, "%d", &height); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的区块高度"})
		return
	}

	// 遍历区块链查找指定高度的区块
	iter := s.blockchain.Iterator()
	var foundBlock *models.Block

	for {
		block := iter.Next()

		if block.Height == height {
			foundBlock = block
			break
		}

		// 如果当前区块高度低于目标高度，表示没有找到
		if block.Height < height {
			break
		}

		// 到达创世区块，退出循环
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	if foundBlock == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到指定高度的区块"})
		return
	}

	c.JSON(http.StatusOK, convertBlockToResponse(foundBlock))
}

// @Summary 获取所有交易
// @Description 返回区块链中的所有交易
// @Tags 交易
// @Accept json
// @Produce json
// @Param limit query int false "限制返回的交易数量"
// @Param offset query int false "跳过前面的交易数量"
// @Success 200 {array} TransactionOutput
// @Router /transactions [get]
func (s *BlockchainServer) getTransactions(c *gin.Context) {
	// 解析分页参数
	limit := 20 // 默认返回20个交易
	offset := 0

	var transactions []TransactionOutput
	currentTx := 0

	iter := s.blockchain.Iterator()
	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if currentTx >= offset {
				transactions = append(transactions, convertTransactionToOutput(tx, block.Timestamp))
			}

			currentTx++

			if limit > 0 && len(transactions) >= limit {
				c.JSON(http.StatusOK, transactions)
				return
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	c.JSON(http.StatusOK, transactions)
}

// @Summary 获取特定交易
// @Description 通过交易ID获取交易详情
// @Tags 交易
// @Accept json
// @Produce json
// @Param id path string true "交易ID"
// @Success 200 {object} TransactionOutput
// @Failure 404 {object} map[string]string
// @Router /transactions/{id} [get]
func (s *BlockchainServer) getTransaction(c *gin.Context) {
	id := c.Param("id")

	// 查找交易
	iter := s.blockchain.Iterator()
	found := false
	var foundTx *models.Transaction
	var blockTime int64

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if fmt.Sprintf("%x", tx.ID) == id {
				found = true
				foundTx = tx
				blockTime = block.Timestamp
				break
			}
		}

		if found {
			break
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "交易未找到"})
		return
	}

	c.JSON(http.StatusOK, convertTransactionToOutput(foundTx, blockTime))
}

// TransactionRequest 创建交易请求结构
type TransactionRequest struct {
	From   string `json:"from" binding:"required"`
	To     string `json:"to" binding:"required"`
	Amount int    `json:"amount" binding:"required,gt=0"`
}

// @Summary 创建新交易
// @Description 创建一个新交易并广播到区块链网络
// @Tags 交易
// @Accept json
// @Produce json
// @Param transaction body TransactionRequest true "交易详情"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /transactions [post]
func (s *BlockchainServer) createTransaction(c *gin.Context) {
	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	// 验证地址
	if !models.ValidateAddress(req.From) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "发送地址格式错误"})
		return
	}

	if !models.ValidateAddress(req.To) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "接收地址格式错误"})
		return
	}

	// 创建交易
	tx := models.NewUTXOTransaction(req.From, req.To, req.Amount, s.blockchain)

	// 挖掘新区块
	newBlock := s.blockchain.MineBlock([]*models.Transaction{tx})

	c.JSON(http.StatusCreated, gin.H{
		"message": "交易已创建",
		"txid":    fmt.Sprintf("%x", tx.ID),
		"block":   fmt.Sprintf("%x", newBlock.Hash),
	})
}

// @Summary 获取所有钱包地址
// @Description 返回系统中所有钱包地址的列表
// @Tags 钱包
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /wallets [get]
func (s *BlockchainServer) getWallets(c *gin.Context) {
	wallets, err := models.NewWallets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取钱包列表失败"})
		return
	}

	addresses := wallets.GetAddresses()
	c.JSON(http.StatusOK, addresses)
}

// @Summary 创建新钱包
// @Description 创建一个新的区块链钱包地址
// @Tags 钱包
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string
// @Router /wallets [post]
func (s *BlockchainServer) createWallet(c *gin.Context) {
	wallets, err := models.NewWallets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建钱包失败"})
		return
	}

	address := wallets.CreateWallet()
	wallets.Save()

	c.JSON(http.StatusCreated, gin.H{"address": address})
}

// @Summary 查询钱包余额
// @Description 获取指定钱包地址的余额
// @Tags 钱包
// @Accept json
// @Produce json
// @Param address path string true "钱包地址"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Router /wallets/{address}/balance [get]
func (s *BlockchainServer) getBalance(c *gin.Context) {
	address := c.Param("address")

	if !models.ValidateAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "地址格式错误"})
		return
	}

	balance := 0
	pubKeyHash := models.GetPubKeyHashFromAddress(address)
	UTXOs := s.blockchain.FindUTXO()

	for _, outputs := range UTXOs {
		for _, out := range outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				balance += out.Value
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// 辅助函数：转换Block为API响应格式
func convertBlockToResponse(block *models.Block) BlockResponse {
	var transactions []TransactionOutput
	for _, tx := range block.Transactions {
		transactions = append(transactions, convertTransactionToOutput(tx, block.Timestamp))
	}

	return BlockResponse{
		Hash:          fmt.Sprintf("%x", block.Hash),
		PrevBlockHash: fmt.Sprintf("%x", block.PrevBlockHash),
		Timestamp:     block.Timestamp,
		Height:        block.Height,
		Nonce:         block.Nonce,
		Transactions:  transactions,
	}
}

// 辅助函数：转换Transaction为API响应格式
func convertTransactionToOutput(tx *models.Transaction, blockTime int64) TransactionOutput {
	var from, to string
	amount := 0

	if tx.IsCoinbase() {
		from = "系统"
		to = fmt.Sprintf("%s", tx.Vout[0].PubKeyHash)
		amount = tx.Vout[0].Value
	} else {
		for _, vin := range tx.Vin {
			from = fmt.Sprintf("%s", vin.PubKey)
			break
		}

		for _, vout := range tx.Vout {
			// 跳过找零输出
			if !bytes.Equal(vout.PubKeyHash, []byte(from)) {
				to = fmt.Sprintf("%s", vout.PubKeyHash)
				amount = vout.Value
				break
			}
		}
	}

	var inputs []string
	for _, vin := range tx.Vin {
		if len(vin.Txid) > 0 {
			inputs = append(inputs, fmt.Sprintf("%x:%d", vin.Txid, vin.Vout))
		}
	}

	return TransactionOutput{
		ID:        fmt.Sprintf("%x", tx.ID),
		From:      from,
		To:        to,
		Amount:    amount,
		Timestamp: blockTime,
		Inputs:    inputs,
	}
}

// 钱包交易历史查询
func (s *BlockchainServer) getWalletTransactions(c *gin.Context) {
	address := c.Param("address")

	if !models.ValidateAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "地址格式错误"})
		return
	}

	var transactions []TransactionOutput
	pubKeyHash := models.GetPubKeyHashFromAddress(address)

	// 遍历区块链
	iter := s.blockchain.Iterator()
	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			isRelevant := false

			// 检查是否为发送方
			if !tx.IsCoinbase() {
				for _, vin := range tx.Vin {
					if bytes.Equal(vin.PubKey, pubKeyHash) {
						isRelevant = true
						break
					}
				}
			}

			// 检查是否为接收方
			if !isRelevant {
				for _, vout := range tx.Vout {
					if vout.IsLockedWithKey(pubKeyHash) {
						isRelevant = true
						break
					}
				}
			}

			// 如果交易与当前地址相关，添加到结果
			if isRelevant {
				transactions = append(transactions, convertTransactionToOutput(tx, block.Timestamp))
			}
		}

		// 如果到达创世区块，退出循环
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	c.JSON(http.StatusOK, transactions)
}
