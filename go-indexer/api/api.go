package api

import (
	"log"
	"net/http"

	"github.com/ichillyzhong/ethereum-defi/go-indexer/db"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures Gin routes
func SetupRouter(dbClient *db.DB) *gin.Engine {
	r := gin.Default()

	// Query Total Value Locked (TVL)
	r.GET("/api/tvl", func(c *gin.Context) {
		tvl, err := dbClient.GetTotalValueLocked()
		if err != nil {
			log.Printf("Error getting TVL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get TVL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tvl": tvl.String(), // Return as string to avoid big number precision issues
		})
	})

	return r
}
