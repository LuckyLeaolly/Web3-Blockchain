package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 简单的全局网络配置
var networkConfig = struct {
	BlockTime    int `json:"blockTime"`
	Difficulty   int `json:"difficulty"`
	MiningReward int `json:"miningReward"`
	TxFee        int `json:"txFee"`
}{
	BlockTime:    5,  // 默认5秒
	Difficulty:   20, // 默认20位
	MiningReward: 50, // 默认50个代币
	TxFee:        1,  // 默认1个代币
}

// 获取网络配置
func (s *BlockchainServer) getNetworkConfig(c *gin.Context) {
	c.JSON(http.StatusOK, networkConfig)
}

// 更新网络配置
func (s *BlockchainServer) updateNetworkConfig(c *gin.Context) {
	var config struct {
		BlockTime    int `json:"blockTime"`
		Difficulty   int `json:"difficulty"`
		MiningReward int `json:"miningReward"`
		TxFee        int `json:"txFee"`
	}

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	// 简单验证
	if config.BlockTime <= 0 || config.Difficulty <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "区块时间和难度必须大于0"})
		return
	}

	// 更新配置
	networkConfig = config

	c.JSON(http.StatusOK, networkConfig)
}
