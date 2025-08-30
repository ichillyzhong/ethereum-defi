package api

import (
	"log"
	"net/http"

	"github.com/ichillyzhong/ethereum-defi/go-indexer/db"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置 Gin 路由
func SetupRouter(dbClient *db.DB) *gin.Engine {
	r := gin.Default()

	// 查询总锁仓价值 (TVL)
	r.GET("/api/tvl", func(c *gin.Context) {
		tvl, err := dbClient.GetTotalValueLocked()
		if err != nil {
			log.Printf("Error getting TVL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get TVL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tvl": tvl.String(), // 返回字符串格式以避免大数精度问题
		})
	})

	return r
}
