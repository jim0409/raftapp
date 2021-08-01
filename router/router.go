package router

import (
	"net/http"
	"raft-app/service"

	"github.com/gin-gonic/gin"
)

// https://stackoverflow.com/questions/29418478/go-gin-framework-cors
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func ApiRouter(r *gin.Engine) {
	authrized := r.Group("/")
	authrized.Use(CORSMiddleware())

	k1 := authrized.Group("/key")
	{
		k1.GET("/:ky", service.GetKey)
		k1.PUT("", service.PutKey)
	}

	n1 := authrized.Group("/node")
	{
		n1.POST("/:nd", service.AddNode)
		n1.DELETE("/:nd", service.DelNode)
	}
}
