package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Middleware de cache (exemplo)
func CacheMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lógica de cache aqui
		c.Next()
	}
}
