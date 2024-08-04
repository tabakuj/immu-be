package handlers

import "github.com/gin-gonic/gin"

func SetupHealth(router *gin.Engine) {
	// now this is a simplifies version of this but in a real life application i would check if we are connected to immudb vault also(typically check all the depenecies)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

}
