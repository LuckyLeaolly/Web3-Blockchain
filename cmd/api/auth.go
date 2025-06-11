package main

import (
	"fmt"
	"net/http"
	"time"

	"personal/web3-blockchain/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// 用于JWT签名的密钥
var jwtSecret = []byte("web3_blockchain_secret_key")

// 简单的内存用户存储
var users = map[string]string{
	"admin": "admin123", // 用户名:密码
}

// JWTAuthMiddleware 创建简单的JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 处理Bearer前缀
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// 解析和验证令牌
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效令牌: " + err.Error()})
			c.Abort()
			return
		}

		// 验证令牌有效性
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 检查令牌是否过期
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌已过期"})
					c.Abort()
					return
				}
			}

			// 将用户信息保存到上下文
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效令牌"})
			c.Abort()
			return
		}
	}
}

// 生成JWT令牌
func generateToken(username string) (string, int64, error) {
	// 使用简单密钥，实际应从配置加载
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	return tokenString, expirationTime.Unix(), err
}

// 登录处理函数
func (s *BlockchainServer) login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	// 简单验证
	password, exists := users[req.Username]
	if !exists || password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成令牌
	token, expires, err := generateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"userId":    1, // 简单实现，固定返回ID 1
		"expiresAt": expires,
	})
}

// 注册处理函数
func (s *BlockchainServer) register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	// 检查用户是否已存在
	if _, exists := users[req.Username]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 保存用户
	users[req.Username] = req.Password

	// 创建钱包
	wallets, _ := models.NewWallets()
	address := wallets.CreateWallet()
	wallets.Save()

	c.JSON(http.StatusCreated, gin.H{
		"id":            len(users),
		"username":      req.Username,
		"email":         req.Email,
		"walletAddress": address,
	})
}
