package routes

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/pastoapp/astroid-api/internal/utils"
)

func KeyGen(c *gin.Context) {
	pub, priv := utils.Keygen()
	c.JSON(200, gin.H{
		"privateKey": priv,
		"publicKey":  pub,
	})
}

type signRequest struct {
	Data string `json:"nonce" binding:"required"`
	Priv string `json:"privateKey" binding:"required"`
}

func Sign(c *gin.Context) {
	var req signRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	signature, err := utils.SignNonce(req.Priv, req.Data)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"signature": base64.StdEncoding.EncodeToString(signature),
	})
}
